ifdef CI
	include .buildkite/Makefile
else
	include .local/Makefile
endif
