def x!(cmd)
	system(cmd) || abort("Could not run #{cmd}")
end

task :default do
	x!("go install ./...")
	exec "pandemic-nerd-hurd start --month jan --funded-events 4"
end