def x!(cmd)
	system(cmd) || abort("Could not run #{cmd}")
end

task :install do
	x!("go install ./...")
end

task :test do
	x!("go test ./...")
end

task :default => :install do
	exec "pandemic-nerd-hurd start --month jan"
end