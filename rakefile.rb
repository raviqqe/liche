TOTAL_COVERAGE_FILE = 'coverage.txt'.freeze # This path is specified by codecov.
BIN_PATH = File.absolute_path 'bin'

task :deps do
  sh 'go get github.com/alecthomas/gometalinter github.com/mattn/goveralls'
  sh 'gometalinter --install'
  sh 'go get -d -t ./...'
  sh 'gem install rake rubocop'
end

task :build do
  sh 'CGO_ENABLED=0 GOOS=linux go build -o bin/liche'
end

task :fast_unit_test do
  sh 'go test ./...'
end

task :unit_test do
  sh "go test -covermode atomic -coverprofile #{TOTAL_COVERAGE_FILE}"
end

task command_test: :build do
  sh 'bundler install'
  sh %W[bundler exec cucumber
        -r examples/aruba.rb
        PATH=#{BIN_PATH}:$PATH
        examples].join ' '
end

task test: %i[unit_test command_test]

task :format do
  sh 'go fix ./...'
  sh 'go fmt ./...'

  Dir.glob '**/*.go' do |file|
    sh "goimports -w #{file}"
  end

  sh 'rubocop -a'
end

task :lint do
  sh 'gometalinter ./...'
  sh 'rubocop'
end

task install: %i[deps test build] do
  sh 'go get ./...'
end

task default: %i[test build]

task :clean do
  sh 'git clean -dfx'
end
