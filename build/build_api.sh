#!/usr/bin/env bash

project_name="/src/conan/api"

api_dir_path=$GOPATH$project_name

cd $api_dir_path


# 遍历api下的所有目录 读取里面的 proto文件
function read_dir(){
  for file in `ls $1` #注意此处这是两个反引号，表示运行系统命令
    do
      if [ -d $1"/"$file ] #注意此处之间一定要加上空格，否则会报错
      then
        read_dir $1$file
      else
        old_path=`pwd`
        next_path=$1"/"
        next_path=${next_path#*/}
        cd $next_path
#        pwd
        protoc -I=$GOPATH/src  -I=. --gofast_out=. *.proto
        # 执行 protoc-go-inject-tag
        cur_file=$file
        cur_file_suf=${cur_file#*.}
        if [ $cur_file_suf = "proto" ]
        then
          cur_file_name=${cur_file%.*}
          protoc-go-inject-tag -input=./${cur_file_name}.pb.go

          echo `pwd`"/"$file" Done"
        fi
        cd $old_path
      fi
    done
}

read_dir ./

#
