#!/bin/bash

OUT_DIR_GENERATED_PROTO_FILES=.
OUT_DIR_GENERATED_3RD_PARTY_PYTHON_FILES="./../.."
PROTOS_PATH=unset
TARGET="gen-core-cpp"

display_usage() {
    available_targets="$(make -f /app/proto_gen.mk list | grep -v  list)"

    echo " -h --help                          Prints out this help message."
    echo "                                    run this container using:"
    echo ""
    echo "                                        \$ docker run -v \$(pwd):\$(pwd) -w \$(pwd) <container_name> -p <dir_containing_protofiles> -d <out_dir> -t <target>"
    echo ""
    echo "                                    Note, all output, input files must be relative to current working directory"
    echo "                                    and cannot be outside current working directory (i.e no ./../../..)"
    echo ""
    echo " -p --proto_path                    Path to the .proto files for which to generate code for relative to the CWD"
    echo " -d --gen_files_out                 Directory where the generated proto files will be outputted. Must be inside current working directory."
    echo "                                    This out directory will be created it if doesn't exist"
    echo " -e --gen_files_out_python_extra    Directory relative to the docker container used to generate python bindings for 3rd party proto files"
    echo "                                    to the directory passed to the container using -w. Must be inside the current working directory"
    echo " -t --target                        List of targets that can be used to generate files for a given language/app"
    echo "                                      Available targets:"
    printf "                                        - %s\n" $available_targets
}

while [[ $# -gt 0 ]]
do
    cmd_arg="$1"

        case $cmd_arg in
            -p|--proto_path)
                PROTOS_PATH=$2
                shift
                shift
            ;;
            -d|--gen_files_out)
                OUT_DIR_GENERATED_PROTO_FILES=$2
                shift
                shift
            ;;
            -e|--gen_files_out_python_extra)
                OUT_DIR_GENERATED_3RD_PARTY_PYTHON_FILES=$2
                shift
                shift
            ;;
            -t|--target)
                TARGET=$2
                shift
                shift
            ;;
            -h|--help)
                display_usage
                exit 0
            ;;
            *)   # unknown option
                echo "Unknown command line option $cmd_arg"
                display_usage
                exit 1
                shift
            ;;
        esac
done

echo "+=================================================================="
echo "| Target                                = $TARGET"
echo "| Path to Live Planet protofiles        = $PROTOS_PATH"
echo "| Out_dir                               = $OUT_DIR_GENERATED_PROTO_FILES"
echo "| Out_dir 3rd party proto (python only) = $OUT_DIR_GENERATED_3RD_PARTY_PYTHON_FILES"
echo "+=================================================================="
pip3 list 2>/dev/null | awk '{ print "| " $0; }'
echo "+=================================================================="

export GOPATH=/go_workspace
export OUT_DIR_GENERATED_PROTO_FILES
export OUT_DIR_GENERATED_3RD_PARTY_PYTHON_FILES
export PROTOS_PATH

make -f /app/proto_gen.mk $TARGET
