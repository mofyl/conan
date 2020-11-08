#!/usr/bin/env bash



project_name="/src/conan/cmd"

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
        go build *.go
        if [ -x $file ]
        then
          ./$file
        fi
#        echo $1$file
#        old_path=`pwd`
#        next_path=$1"/"
#        next_path=${next_path#*/}
#        cd $next_path
#        echo $next_path
##        pwd
##        protoc -I=. -I=/Users/mo/go/src/ --gofast_out=. *.proto
#        # 执行 protoc-go-inject-tag
##
#        cd $old_path
      fi
    done
}

read_dir ./


