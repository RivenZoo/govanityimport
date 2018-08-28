# 通过下面的变量指定需要测试的包和不需要的包
# 目录路径
TEST_DIRS += ./controllers
# 只支持包路径
EXCLUDE_PKG :=

TEST_GEN_CMD := gotests
TEST_GEN_FLAGS += -all -w
TEST_RUN_FLAGS += -race -cover -v

EXCLUDE_FILE_FLAG = $(foreach n,$(EXCLUDE_PKG),! -path '$(n)*')

TEST_SRC_FILES = $(shell cd $(PROJ_GOPATH); find $(TEST_DIRS) -type f -name '*.go' ! -name '*_test.go' $(EXCLUDE_FILE_FLAG))
TEST_CASE_PKG = $(sort $(dir $(TEST_SRC_FILES)))
TEST_FILES = $(subst .go,_test.go,$(TEST_SRC_FILES))

test_case: 
	@echo "generate test case for packages:" 
	@echo "$(TEST_CASE_PKG)"
	cd $(PROJ_GOPATH); $(TEST_GEN_CMD) $(TEST_GEN_FLAGS) $(TEST_CASE_PKG)

test_run:
	cd $(PROJ_GOPATH); go test $(TEST_RUN_FLAGS) $(TEST_CASE_PKG)


