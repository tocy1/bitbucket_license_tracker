pipeline {
    agent {
            kubernetes {
            inheritFrom 'docker'
            yamlMergeStrategy override()
            containerTemplate(name: 'jnlp',
             image: 'ephemeral-docker-virtual.artifactory.swisscom.com/jenkins-slave-internal-dind:2.0.3')
            }
            }
    options { buildDiscarder(logRotator(numToKeepStr: '6')) }
    environment {
        ARTIFACTORY_URL='dos-docker.artifactory.swisscom.com'
        IMAGE_NAME='prom-bb-licensce'
        IMAGE_TAG='1.0.1'
        SNYK_TOKEN=credentials('dsoe-snyk-test-sa-secrettext-token')
        ARTIFACTORY=credentials('saas-ops-sa-secret')
        SNYK_SEVERITY_THRESHOLD='high'
    }
    stages {
        stage('Build Image') {
            steps {
              echo 'Building docker image'
              sh "docker build --no-cache -t ${ARTIFACTORY_URL}/${IMAGE_NAME}:${IMAGE_TAG} ."
            }
        }
        stage('Snyk Container') {
            steps {
                echo "Scanning container image for vulnerabilities"
                sh "snyk container  test --severity-threshold=${SNYK_SEVERITY_THRESHOLD} ${ARTIFACTORY_URL}/${IMAGE_NAME}:${IMAGE_TAG}"
            }
        }
        stage('Push Image') {
            steps {
              sh "docker login -u ${ARTIFACTORY_USR} -p ${ARTIFACTORY_PSW} ${ARTIFACTORY_URL}"
              echo "Pushing docker image ${ARTIFACTORY_URL}/${IMAGE_NAME}:${IMAGE_TAG}"
              sh "docker push ${ARTIFACTORY_URL}/${IMAGE_NAME}:${IMAGE_TAG}"
            }
        }
    }
    post {
        always {
            cleanWs()
        }
    }
}
