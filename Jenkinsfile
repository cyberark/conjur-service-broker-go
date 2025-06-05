#!/usr/bin/env groovy

@Library(['product-pipelines-shared-library', 'conjur-enterprise-sharedlib']) _

def productName = 'Conjur Service Broker Go'
def productTypeName = 'Conjur Enterprise'

// Automated release, promotion and dependencies
properties([
  // Include the automated release parameters for the build
  release.addParams(),
  // Dependencies of the project that should trigger builds
  dependencies([
    'conjur-enterprise/conjur-api-go',
  ])
])

// Performs release promotion.  No other stages will be run
if (params.MODE == "PROMOTE") {
  release.promote(params.VERSION_TO_PROMOTE) { infrapool, sourceVersion, targetVersion, assetDirectory ->
    // Any assets from sourceVersion Github release are available in assetDirectory
    // Any version number updates from sourceVersion to targetVersion occur here
    // Any publishing of targetVersion artifacts occur here
    // Anything added to assetDirectory will be attached to the Github Release

    def sourceTag = infrapool.agentSh(
        returnStdout: true,
        script: "source ./scripts/build_utils.sh && echo \"${sourceVersion}-\$(git_commit)\""
    ).trim()

    env.INFRAPOOL_PRODUCT_NAME = "${productName}"
    env.INFRAPOOL_DD_PRODUCT_TYPE_NAME = "${productTypeName}"

    // Scan the image before promoting
    runSecurityScans(infrapool,
      image: "registry.tld/conjur-service-broker:${sourceTag}",
      buildMode: params.MODE,
      branch: env.BRANCH_NAME,
      arch: 'linux/amd64')

    // Promote source version to target version.

    // NOTE: the use of --pull to ensure source images are pulled from internal registry
    infrapool.agentSh "./scripts/publish_container_images.sh --promote --source ${sourceTag} --target ${targetVersion} --pull"
  }
  release.copyEnterpriseRelease(params.VERSION_TO_PROMOTE)
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

    // Values to direct scan results to the right place in DefectDojo
    INFRAPOOL_PRODUCT_NAME = "${productName}"
    INFRAPOOL_DD_PRODUCT_TYPE_NAME = "${productTypeName}"
  }

  triggers {
    cron(getDailyCronString())
    parameterizedCron(getWeeklyCronString("H(1-5)","%MODE=RELEASE"))
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

    stage('Scan for internal URLs') {
      steps {
        script {
          detectInternalUrls()
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

    // Generates a VERSION file based on the current build number and latest version in CHANGELOG.md
    stage('Validate changelog and set version') {
      steps {
        script {
          updateVersion(infrapool, "CHANGELOG.md", "${BUILD_NUMBER}")
        }
      }
    }

    stage('Get latest upstream dependencies') {
      steps {
        script {
          updatePrivateGoDependencies("${WORKSPACE}/go.mod")
          // Copy the vendor directory onto infrapool
          infrapool.agentPut from: "vendor", to: "${WORKSPACE}"
          infrapool.agentPut from: "go.*", to: "${WORKSPACE}"
          // Add GOMODCACHE directory to infrapool allowing automated release to generate SBOMs
          infrapool.agentPut from: "/root/go", to: "/var/lib/jenkins/"
        }
      }
    }

    stage('Build while unit testing') {
      parallel {
        stage('Run unit tests') {
          steps {
            script {
              infrapool.agentSh './scripts/test_in_docker.sh --skip-gomod-download'
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
              infrapool.agentSh "./scripts/build_release.sh --skip=validate --clean"

              // Build container images
              infrapool.agentSh "./scripts/build_container_images.sh"

              // Archive release artifacts
              infrapool.agentArchiveArtifacts artifacts: 'dist/goreleaser/'
            }
          }
        }
      }
    }

    stage('Integration tests') {
      steps {
        script {
          infrapool.agentSh './scripts/test_integration.sh --skip-gomod-download'
        }
      }
    }

    stage("Push container images to internal registry") {
      steps {
        script {
          // Publish container images to internal registry
          infrapool.agentSh './scripts/publish_container_images.sh --internal'
        }
      }
    }

    stage('End to end test while scanning') {
      parallel {
        // stage('End-to-End testing') {
        //   steps {
        //     script {
        //       allocateTas(infrapool, 'isv_ci_tas_srt_5_0')
        //       infrapool.agentSh './test/e2e/test.sh'
        //       infrapool.agentStash name: 'e2e-test-results', includes: 'test/e2e/reports/junit.xml'
        //     }
        //   }

        //   post {
        //     always {
        //       destroyTas(infrapool)
        //       unstash 'e2e-test-results'
        //       junit 'test/e2e/reports/junit.xml'
        //     }
        //   }
        // }

        stage("Scan container images for fixable issues") {
          steps {
             runSecurityScans(infrapool,
                image: "registry.tld/${containerImageWithTag()}",
                buildMode: params.MODE,
                branch: env.BRANCH_NAME,
                arch: 'linux/amd64')
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
            infrapool.agentSh "cp -r dist/goreleaser/*.zip dist/goreleaser/SHA256SUMS.txt ${assetDirectory}"

            // Create Go application SBOM using the go.mod version for the golang container image
            infrapool.agentSh """export PATH="${toolsDirectory}/bin:${PATH}" && go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --main "cmd/conjur_service_broker/" --output "${billOfMaterialsDirectory}/go-app-bom.json" """
            // Create Go module SBOM
            infrapool.agentSh """export PATH="${toolsDirectory}/bin:${PATH}" && go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --output "${billOfMaterialsDirectory}/go-mod-bom.json" """
          })
        }
      }
    }
  }
  post {
    always {
      releaseInfraPoolAgent(".infrapool/release_agents")

      // Resolve ownership issue before running infra post hook
      sh 'git config --global --add safe.directory ${PWD}'
      infraPostHook()
    }
  }
}

def containerImageWithTag() {
  infrapool.agentSh(
    returnStdout: true,
    script: 'source ./scripts/build_utils.sh && echo "conjur-service-broker:$(project_version_with_commit)"'
  )
}
