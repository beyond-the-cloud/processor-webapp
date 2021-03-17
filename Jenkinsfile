pipeline {
    agent any
    triggers {
        githubPush()
    }
    stages {
        stage('Build') {
            steps {
                echo 'Building the application...'
                build job: 'BuildProcessorWebapp'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying the application...'
                build job: 'BuildProcessorWebapp'
            }
        }
    }
}