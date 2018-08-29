.PHONY: pytest

PYTEST = $(shell which pytest)
TEST_DIR = .
TEST_SRC = $(wildcard $(TEST_DIR)/*.py)
TEST_RESULT = ./result.log

ifeq ($(PYTEST),)
$(error No pytest found!)
endif

pytest: $(TEST_RESULT)

$(TEST_RESULT): $(TEST_SRC)
	$(PYTEST) --capture=no $? > $@

