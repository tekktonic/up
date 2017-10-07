
set SITE=localhost:8080

check()
{
    if [ $? = 0 ]; then
        echo $@;
        exit 1;
    fi;
}


die()
{
    echo $?
    exit 1
}

noauthpost()
{
    return curl  -w "%{http_version}\n" $1 -d $2=$3 2>/dev/null | tail -n 1 | grep 401
}

badauthpost() {
    return curl -H "X-Up-Auth: $2" -w "%{http_version}\n" $1 -d $3=$4 2>/dev/null | tail -n 1 | grep 401
}


noauthget()
{
    return curl  -w "%{http_version}\n" $1 2>/dev/null | tail -n 1 | grep 401
}

badauthget() {
    return curl -H "X-Up-Auth: $2" -w "%{http_version}\n" $1 2>/dev/null | tail -n 1 | grep 401
}


setuptest()
{
    mkdir testing_internal && cd testing_internal
    if [ `pgrep up` ]; then
        echo "Killing Up so we can test"
        pkill up
    fi

    cp ../../config.json . || die "Unable to copy config file from `pwd`../../"
    ../../up &
    sleep 5;
}

teardowntest()
{
    if [ `pgrep up` ]; then
        echo "Cleaning up up"
        pkill up
    fi

    cd ..
    rm -rf testing_internal || die "Somehow got paths confused"
}
