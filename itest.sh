#!/bin/sh

bin=kubecfg
target=tests/out.yaml
if test -f "$bin"; then
    rm $target
    touch $target
    if ! ./$bin import tests/resources/test.yaml -t $target -n context; then 
        echo "import test.yaml to $target failed"
        exit 1
    fi

    if ! ./$bin switch context-1 -t $target; then 
        echo "switch context failed"
        exit 1
    fi

    if ! ./$bin list -t $target; then 
        echo "list all contexts failed"
        exit 1
    fi

    if ! ./$bin view context -t $target; then 
        echo "view 'context' failed"
        exit 1
    fi

    if ! ./$bin view -t $target; then 
        echo "view active context failed"
        exit 1
    fi

    if ! ./$bin delete context -t $target; then
        echo "delete 'context' failed"
        exit 1
    fi

    exportTarget=tests/export.yaml
    rm $exportTarget    
    if ./$bin export context-1 -s $target -t $exportTarget; then         
        if ! cmp -s "$exportTarget" "tests/expected/export.yaml"; then
            cmp -s "$exportTarget" "tests/expected/export.yaml"
            exit 1
        fi
    else
        echo "export 'context-1' failed"
        exit 1
    fi

    if ./$bin rename context-1 foo -t $target; then
        if ! cmp -s "$target" "tests/expected/out.yaml"; then
            cmp -s "$target" "tests/expected/out.yaml"
            exit 1
        fi
    else
        echo "rename 'context-1' failed"
        exit 1
    fi
else
    echo "kubecfg binary not found."
fi