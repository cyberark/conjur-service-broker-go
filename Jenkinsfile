#!/usr/bin/env groovy

@Library("product-pipelines-shared-library") _

// Automated release, promotion and dependencies
properties([
  // Include the automated release parameters for the build
  release.addParams(),
  // Dependencies of the project that should trigger builds
  dependencies([
    'cyberark/conjur-api-go',
  ])
])

// Performs release promotion.  No other stages will be run
if (params.MODE == "PROMOTE") {
  release.promote(params.VERSION_TO_PROMOTE) { infrapool, sourceVersion, targetVersion, assetDirectory ->
    // Any assets from sourceVersion Github release are available in assetDirectory
    // Any version number updates from sourceVersion to targetVersion occur here
    // Any publishing of targetVersion artifacts occur here
    // Anything added to assetDirectory will be attached to the Github Release

    // Promote source version to target version.

    // NOTE: the use of --pull to ensure source images are pulled from internal registry
    infrapool.agentSh "source ./scripts/build_utils.sh && ./scripts/publish_container_images.sh --promote --source ${sourceVersion}-\$(git_commit) --target ${targetVersion} --pull"
  }
  return
}

pipeline {
  agent { label 'conjur-enterprise-common-agent' }

  options {
    timestamps()
    buildDiscarder(logRotator(numToKeepStr: '30'))
    timeout(time: 30, unit: 'MINUTES')
  }

  environment {
    // Sets the MODE to the specified or autocalculated value as appropriate
    MODE = release.canonicalizeMode()
  }

  triggers {
    cron(getDailyCronString())
  }

  stages {
    // Aborts any builds triggered by another project that wouldn't include any changes
    stage ("Skip build if triggering job didn't create a release") {
      when {
        expression {
          MODE == "SKIP"
        }
      }
      steps {
        script {
          currentBuild.result = 'ABORTED'
          error("Aborting build because this build was triggered from upstream, but no release was built")
        }
      }
    }

    stage('Get InfraPool ExecutorV2 Agent(s)') {
      steps{
        script {
          // Request ExecutorV2 agents for 1 hour(s)
          INFRAPOOL_EXECUTORV2_AGENTS = getInfraPoolAgent(type: "ExecutorV2", quantity: 1, duration: 1)
          INFRAPOOL_EXECUTORV2_AGENT_0 = INFRAPOOL_EXECUTORV2_AGENTS[0]
          infrapool = infraPoolConnect(INFRAPOOL_EXECUTORV2_AGENT_0, {})
        }
      }
    }

    stage('Parse Changelog') {
      steps {
        script {
          infrapool.agentSh './scripts/parse-changelog.sh'
        }
      }
    }

    // Generates a VERSION file based on the current build number and latest version in CHANGELOG.md
    stage('Validate changelog and set version') {
      steps {
        script {
          updateVersion(infrapool, "CHANGELOG.md", "${BUILD_NUMBER}")
        }
      }
    }

    stage('Build while unit testing') {
      parallel {
        stage('Run unit tests') {
          steps {
            script {
              infrapool.agentSh './scripts/test_in_docker.sh'
              infrapool.agentStash name: 'test-results', includes: 'coverage/*'
            }
          }

          post {
            always {
              unstash 'test-results'
              junit 'coverage/junit.xml'
              cobertura(
                autoUpdateHealth: false,
                autoUpdateStability: false,
                coberturaReportFile: 'coverage/cobertura.xml',
                conditionalCoverageTargets: '70, 0, 0',
                failUnhealthy: false,
                failUnstable: false,
                maxNumberOfBuilds: 0,
                lineCoverageTargets: '70, 0, 0',
                methodCoverageTargets: '70, 0, 0',
                onlyStable: false,
                sourceEncoding: 'ASCII',
                zoomCoverageChart: false
              )

              // Don't fail builds if we can't upload coverage information to Codacy
              catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                retry (2) {
                  codacy action: 'reportCoverage', language: 'Go', filePath: 'coverage/all_no_gen', extraArgs: '--force-coverage-parser go'
                }
              }
            }
          }
        }

        stage('Build release artifacts') {
          steps {
            script {
              // Create release artifacts without releasing to Github
              infrapool.agentSh "./scripts/build_release.sh --skip-validate --clean"

              // Build container images
              infrapool.agentSh "./scripts/build_container_images.sh"

              // Archive release artifacts
              infrapool.agentArchiveArtifacts artifacts: 'dist/goreleaser/'
            }
          }
        }
      }
    }

    stage('End to end test while scanning') {
      parallel {
        stage('End-to-End testing') {
          steps {
            script {
              allocateTas(infrapool, 'isv_ci_tas_srt_3_0')
              infrapool.agentSh './test/e2e/test.sh'
              infrapool.agentStash name: 'e2e-test-results', includes: 'test/e2e/reports/junit.xml'
            }
          }

          post {
            always {
              destroyTas(infrapool)
              unstash 'e2e-test-results'
              junit 'test/e2e/reports/junit.xml'
            }
          }
        }

        stage("Scan container images for fixable issues") {
          steps {
             scanAndReport(infrapool, "${containerImageWithTag()}", "HIGH", false)
          }
        }

        stage("Scan container images for total issues") {
          steps {
            scanAndReport(infrapool, "${containerImageWithTag()}", "HIGH", false)
          }
        }
      }
    }

    stage('Release') {
      when {
        expression {
          MODE == "RELEASE"
        }
      }

      steps {
        script {
          release(infrapool, { billOfMaterialsDirectory, assetDirectory, toolsDirectory ->
            // Publish release artifacts to all the appropriate locations
            // Copy any artifacts to assetDirectory to attach them to the Github release

            // Copy assets to be published in Github release.
            infrapool.agentSh "./scripts/copy_release_artifacts.sh ${assetDirectory}"

            // Create Go application SBOM using the go.mod version for the golang container image
            infrapool.agentSh """go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --main "cmd/conjur_service_broker/" --output "${billOfMaterialsDirectory}/go-app-bom.json" """
            // Create Go module SBOM
            infrapool.agentSh """go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --output "${billOfMaterialsDirectory}/go-mod-bom.json" """

            // Publish container images to internal registry
            infrapool.agentSh './scripts/publish_container_images.sh --internal'
          })
        }
      }
    }
  }
  post {
    always {
      releaseInfraPoolAgent(".infrapool/release_agents")
    }
  }
}

def containerImageWithTag() {
  infrapool.agentSh(
    returnStdout: true,
    script: 'source ./scripts/build_utils.sh && echo "conjur-service-broker:$(project_version_with_commit)"'
  )
}
