#!groovy
dockerhubCred = [
    $class: 'UsernamePasswordMultiBinding',
    usernameVariable: 'DOCKERHUB_USERNAME',
    passwordVariable: 'DOCKERHUB_PASSWORD',
    credentialsId: 'dshish-dockerhub-up',
]

pipeline {
    parameters {
        string(name: 'VERSION', defaultValue: '', description: 'Version of the application')
    }

    agent { label 'pod && linux && amd64' }
    stages {
        stage('Test') {
            steps {
                sh '''
                    docker run --rm -v $(pwd):/app -w /app golang:1.19.0 go test -v ./user
                '''
                 sh '''
                    docker run --rm -v $(pwd):/app -w /app golang:1.19.0 go test -v ./message
                '''
            }
        }

        stage('Build') {
            steps {
                sh "./scripts/build.sh ${params.VERSION}"
            }
        }

        stage('Push') {
            steps {
                withCredentials([dockerhubCred]) {
                    sh """
                        docker login -u \$DOCKERHUB_USERNAME -p \$DOCKERHUB_PASSWORD
                        docker push stwongmirantis/messaging-server:${params.VERSION}
                    """
                }
            }
        }
    }

    post {
        always {
            sh("docker rmi stwongmirantis/messaging-server:${params.VERSION}")
        }
    }
}