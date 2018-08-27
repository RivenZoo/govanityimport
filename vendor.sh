#!/bin/sh

GOVENDOR=$(which govendor)

# internal function
function log {
	echo $(date +"%Y%m%d %H:%M:%S") " $@"
}

function is_git_repo_url {
	param=$1
	if [[ "${param:0:4}" = "git:" ]]; then
		return 1
	fi
	if [[ "${param:0:7}" = "http://" ]]; then
		return 1
	fi
	if [[ "${param:0:8}" = "https://" ]]; then
		return 1
	fi
	return 0
}

function add_by_git_repo_url {
	repo=$1
	spec=$2

	repo_name=${repo##*/}
	repo_name=${repo_name%%.*}

	gopath=${GOPATH%%:*}
	local_repo=${gopath}/src/${repo_name}

	if [[ -d ${local_repo} ]]; then
		log "${local_repo} already exist"
		if [[ -n "${spec}" ]]; then
			(cd ${local_repo} && git fetch --all --tags --prune && git checkout ${spec})
		fi
	else
		if [[ -n "${spec}" ]]; then
			git clone --branch ${spec} ${repo} ${local_repo} && git checkout ${spec}
		else
			git clone ${repo} ${local_repo} && git checkout master
		fi
		
	fi
	${GOVENDOR} add ${repo_name}
}

function add_by_module_path {
	repo=$1
	spec=$2

	go get -u ${repo}

	fetch_url=${repo}
	if [[ -n "${spec}" ]]; then
		fetch_url=${fetch_url}@${spec}
	fi
	${GOVENDOR} add ${fetch_url}
}

function add {
	repo=$1
	spec=$2
	is_git_repo_url ${repo}
	is_repo_url=$?

	if [[ ${is_repo_url} -eq 1 ]]; then
		add_by_git_repo_url ${repo} ${spec}
	else
		add_by_module_path ${repo} ${spec}
	fi
}

function update {
	repo=$1
	spec=$2
	is_git_repo_url ${repo}
	is_repo_url=$?

	if [[ ${is_repo_url} -eq 1 ]]; then
		log "need module_import_path"
		exit 0
	fi
	
	gopath=${GOPATH%%:*}
	local_repo=${gopath}/src/${repo}

	fetch_url=${repo}
	update_cmd=cd ${local_repo} && git checkout master && git pull --rebase origin master
	if [[ -n "${spec}" ]]; then
		fetch_url=${fetch_url}@${spec}
		update_cmd=cd ${local_repo} && git fetch --all --tags --prune && git checkout ${spec}
	fi
	(${update_cmd})

	${GOVENDOR} update ${fetch_url}
}

function rm {
	repo=$1
	${GOVENDOR} remove ${repo}
}

function init {
	${GOVENDOR} init
}

function usage {
	cat <<EOF
# vendor.sh init

# vendor.sh add module_url(module path | git url) [git_spec(tags | branch | commit)]

# vendor.sh update module_import_path [git_spec(tags | branch | commit)]

# vendor.sh rm module_import_path
EOF
}

if [[ -z "${GOVENDOR}" ]]; then
	log "no govendor install"
	exit 0
fi

case $1 in
	init) init;;
	add) add $2 $3;;
	update)	update $2 $3;;
	rm)	rm $2;;
	help) usage;;
	*) usage ;;
esac