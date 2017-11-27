set SITE=localhost:9090

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

authpost() {
    return curl -H "X-Up-Auth: $2" -w "%{http_version}\n" $1 -d $3=$4 2>/dev/null | tail -n 1 | grep 200
}

noauthget()
{
    return curl  -w "%{http_version}\n" $1 2>/dev/null | tail -n 1 | grep 401
}

badauthget() {
    return curl -H "X-Up-Auth: $2" -w "%{http_version}\n" $1 2>/dev/null | tail -n 1 | grep 401
}

authget() {
    return curl -H "X-Up-Auth: $2" -w "%{http_version}\n" $1 2>/dev/null | tail -n 1 | grep 200
}

authgetvalue() {
    curl -q -H "X-Up-Auth: $1"  $2 # 2>/dev/null
}

setuptest()
{
    mkdir testing_internal && cd testing_internal || die "Testing internal already exists, something go wrong before? Look into that."
    if [ `pgrep up` ]; then
        echo "Killing Up so we can test"
        pkill up
    fi

    echo "Copying config file in..."
    CONFIGFILE="../../test.config.json"
    if [ ! -e $CONFIGFILE ]; then echo "No test config file"; CONFIGFILE="../../config.json"; fi
    echo $CONFIGFILE
    cp $CONFIGFILE config.json || die "Unable to copy config file from `pwd`../../"
    echo "Generating a fresh db..."
    ../../scripts/gendb.sh test.up.db
    
    ../../up &
    echo "Up's PID is " `pgrep up`
    sleep 5;
}

setuptest2()
{
    mkdir testing_internal && cd testing_internal || die "Testing internal already exists, something go wrong before? Look into that."
    if [ `pgrep up` ]; then
        echo "Killing Up so we can test"
        pkill up
    fi

    echo "Copying config file in..."
    CONFIGFILE="../../test2.config.json"
    if [ ! -e $CONFIGFILE ]; then echo "No test config file"; CONFIGFILE="../../config.json"; fi
    echo $CONFIGFILE
    cp $CONFIGFILE config.json || die "Unable to copy config file from `pwd`../../"
    echo "Generating a fresh db..."
    ../../scripts/gendb.sh test2.up.db
    
    ../../up test2.config.json &
    echo "Up's PID is " `pgrep up`
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
