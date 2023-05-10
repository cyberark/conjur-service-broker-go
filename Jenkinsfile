#!/usr/bin/env groovy

// This is a template Jenkinsfile for builds and the automated release project

// Automated release, promotion and dependencies
properties([
  // Include the automated release parameters for the build
  release.addParams(),
  // Dependencies of the project that should trigger builds
  dependencies(['cyberark/conjur-api-go'])
])

// Performs release promotion.  No other stages will be run
if (params.MODE == "PROMOTE") {
  release.promote(params.VERSION_TO_PROMOTE) { sourceVersion, targetVersion, assetDirectory ->
    // Any assets from sourceVersion Github release are available in assetDirectory
    // Any version number updates from sourceVersion to targetVersion occur here
    // Any publishing of targetVersion artifacts occur here
    // Anything added to assetDirectory will be attached to the Github Release
  }
  return
}

pipeline {
  agent { label 'conjur-enterprise-common-agent' }

  options {
    timestamps()
    buildDiscarder(logRotator(numToKeepStr: '30'))
  }

  triggers {
    cron(getDailyCronString())
  }

  environment {
    // Sets the MODE to the specified or autocalculated value as appropriate
    MODE = release.canonicalizeMode()
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
        }
      }
    }

    // Generates a VERSION file based on the current build number and latest version in CHANGELOG.md
    stage('Validate Changelog and set version') {
      steps {
        script {
          infraPoolConnect(INFRAPOOL_EXECUTORV2_AGENT_0) { infrapool ->
            updateVersion(infrapool, "CHANGELOG.md", "${BUILD_NUMBER}")
          }
        }
      }
    }

//     stage('Get latest upstream dependencies') {
//       steps {
//         updateGoDependencies("${WORKSPACE}/go.mod")
//       }
//     }

    stage('Build and Unit tests') {
      parallel {
        stage('Build binary') {
          steps {
            script {
              infraPoolConnect(INFRAPOOL_EXECUTORV2_AGENT_0) { infrapool ->
                infrapool.agentSh './scripts/build'
              }
            }
          }
        }

        stage('Unit tests') {
          steps {
            script {
              infraPoolConnect(INFRAPOOL_EXECUTORV2_AGENT_0) { infrapool ->
                infrapool.agentSh './scripts/test_in_docker'
                infrapool.agentStash name: 'test-results', includes: 'coverage/*.xml'
              }
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
               lineCoverageTargets: '70, 0, 0',
               maxNumberOfBuilds: 0,
               methodCoverageTargets: '70, 0, 0',
               onlyStable: false,
               sourceEncoding: 'ASCII',
               zoomCoverageChart: false
             )
            }
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
          infraPoolConnect(INFRAPOOL_EXECUTORV2_AGENT_0) { infrapool ->
            release(infrapool) { bomDirectory, assetDirectory ->
            }
          }
        }
      }
    }
  }

  post {
    always {
//       cleanupAndNotify(currentBuild.currentResult)
      releaseInfraPoolAgent(".infrapool/release_agents")
    }
  }
}
