#!/usr/bin/env zsh

export PATH="$PATH:./build/bin"

if ! whence -p refgen &> /dev/null; then
    echo "make build first"
    exit 1
fi


refgen=./build/bin/refgen

MAX_USER_ID=$(psql -qtAXc 'SELECT max(id) + 1 from user_;')

USER_ID=$(
psql -qtAX <<END
INSERT INTO user_
    (ref_id, email, name, pwhash)
VALUES
    (decode('$(refgen generate -b hex -t 1)', 'hex'),
     'user${MAX_USER_ID}@example.com',
     'test-user${MAX_USER_ID}',
     decode('246172676f6e32696424763d3139246d3d36353533362c743d332c703d32242f4f6159544139306b35686759316e36746a7a4d4751244471323959354753484f476a6c4f664b6864475943494f697554344b2f374f757177746e74715278774367', 'hex'))
    RETURNING id
    ;
END
)

psql <<END
INSERT INTO event_ 
    (user_id, ref_id, name, description, start_time)
VALUES 
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 01', 'event 01 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 02', 'event 02 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 03', 'event 03 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 04', 'event 04 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 05', 'event 05 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 06', 'event 06 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 07', 'event 07 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 08', 'event 08 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 09', 'event 09 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 10', 'event 10 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 11', 'event 11 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 12', 'event 12 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 13', 'event 13 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 14', 'event 14 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 15', 'event 15 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 16', 'event 16 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 17', 'event 17 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 18', 'event 18 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 19', 'event 19 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 20', 'event 20 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refgen generate -b hex -t 2)', 'hex'), 'event 21', 'event 21 description', CURRENT_TIMESTAMP)
    ;
END

EVENT_ID=$(psql -qtAXc 'SELECT id from event_ limit 1')

psql <<END
INSERT INTO event_item_ 
    (ref_id, event_id, description)
VALUES 
    (decode('$(refgen generate -b hex -t 3)', 'hex'), ${EVENT_ID}, 'event item 1 description'),
    (decode('$(refgen generate -b hex -t 3)', 'hex'), ${EVENT_ID}, 'event item 2 description'),
    (decode('$(refgen generate -b hex -t 3)', 'hex'), ${EVENT_ID}, 'event item 3 description'),
    (decode('$(refgen generate -b hex -t 3)', 'hex'), ${EVENT_ID}, 'event item 4 description'),
    (decode('$(refgen generate -b hex -t 3)', 'hex'), ${EVENT_ID}, 'event item 5 description')
    ;
END

EVENT_ITEM_ID=$(psql -qtAXc 'SELECT id from event_item_ limit 1')
psql <<END
INSERT INTO earmark_ 
    (ref_id, event_item_id, user_id, note)
VALUES 
    (decode('$(refgen generate -b hex -t 4)', 'hex'), ${EVENT_ITEM_ID}, 1, 'i love pickles!')
    ;
END
