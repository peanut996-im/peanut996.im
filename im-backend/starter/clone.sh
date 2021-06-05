#!/bin/bash
# author: peanut996
# date: 2021.3.31
# des: 克隆本项目并重命名

# 通过此脚本直接将 starter 重命名为自己的服务工程
repository=$1
currepository=${PWD##*/}
oldrepository=${PWD##*/}
if [ "$repository" = "" ] ;then
    echo "没有输入新的項目名"
    exit
fi

# 拷贝的是Demo代码，需要直接改动
if [ "$repository" = "$oldrepository" ] ;then
    oldrepository="svrdemo"
fi
echo "克隆 $oldrepository -> $repository"

check=`ls ../ | grep $repository`
if [ "$check" == "" ] ;then
    echo "项目 $repository 目录不存在"
    cd ../
    echo $(pwd)"/"$repository
    mkdir $repository
    # exit
fi

function rename_go() {
	files=`ls`
    for file in ${files[@]}; do
		if [ -d "$(pwd)/$file" ] ;then
            cd ./$file
            rename_go
            cd ../
        else
            if [[ $file == *.go ]] ;then
                if [[ `cat $(pwd)/$file | grep "\"$oldrepository/"` != "" ]]; then
                    rm -rf $file.tmp
                    sed "s/\"$oldrepository\//\"$repository\//g" $file >> $file.tmp
                    rm -rf $file
                    mv $file.tmp $file
                fi

                if [[ `cat $(pwd)/$file | grep "pb$oldrepository"` != "" ]]; then
                    rm -rf $file.tmp
                    sed "s/pb$oldrepository/pb$repository/g" $file >> $file.tmp
                    rm -rf $file
                    mv $file.tmp $file
                fi
            fi
        fi
	done
}

# step1 复制所有文件
if [ "$repository" != "$currepository" ] ;then
    cd $currepository
    srcs=`ls`
    for src in ${srcs[@]}; do
        if [ "$src" != ".git" ] ;then
            cp -R ./$src ../$repository/$src
        fi
    done
    cd ../$repository
fi

# step2 重命名import
cd ./src
    rename_go
cd ../

# step3 go验证一下
sh go.sh