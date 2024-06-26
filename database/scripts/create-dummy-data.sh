#!/usr/bin/env zsh

set -e
export PATH="$PATH:./build/bin"

if ! whence -p refidtool &> /dev/null; then
    echo "make build first"
    exit 1
fi

NEXT_USER_ID=$(psql -qtAXc 'SELECT max(id) + 1 from user_;')
if [ -z "$NEXT_USER_ID" ]; then
    NEXT_USER_ID=1
fi
echo ">> adding user with id=${NEXT_USER_ID}"

USER_ID=$(
psql -qtAX <<END
INSERT INTO user_
    (ref_id, email, name, pwhash)
VALUES
    (decode('$(refidtool generate -o hex -t 1)', 'hex'),
     'user${NEXT_USER_ID}@example.com',
     'test-user${NEXT_USER_ID}',
     decode('246172676f6e32696424763d3139246d3d36353533362c743d332c703d32242f4f6159544139306b35686759316e36746a7a4d4751244471323959354753484f476a6c4f664b6864475943494f697554344b2f374f757177746e74715278774367', 'hex'))
    RETURNING id
    ;
END
)

NEXT_USER2_ID=$(psql -qtAXc 'SELECT max(id) + 1 from user_;')
echo ">> adding a second user with id=${NEXT_USER2_ID}"

USER_ID2=$(
psql -qtAX <<END
INSERT INTO user_
    (ref_id, email, name, pwhash)
VALUES
    (decode('$(refidtool generate -o hex -t 1)', 'hex'),
     'user${NEXT_USER2_ID}@example.com',
     'test-user${NEXT_USER2_ID}',
     decode('246172676f6e32696424763d3139246d3d36353533362c743d332c703d32242f4f6159544139306b35686759316e36746a7a4d4751244471323959354753484f476a6c4f664b6864475943494f697554344b2f374f757177746e74715278774367', 'hex'))
    RETURNING id
    ;
END
)


echo ">> creating some events for user id=${USER_ID}"
psql -qtAX <<END
INSERT INTO event_ 
    (user_id, ref_id, name, description, start_time)
VALUES 
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 01', 'event 01 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 02', 'event 02 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 03', 'event 03 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 04', 'event 04 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 05', 'event 05 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 06', 'event 06 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 07', 'event 07 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 08', 'event 08 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 09', 'event 09 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 10', 'event 10 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 11', 'event 11 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 12', 'event 12 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 13', 'event 13 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 14', 'event 14 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 15', 'event 15 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 16', 'event 16 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 17', 'event 17 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 18', 'event 18 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 19', 'event 19 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 20', 'event 20 description', CURRENT_TIMESTAMP),
    (${USER_ID}, decode('$(refidtool generate -o hex -t 2)', 'hex'), 'event 21', 'event 21 description', CURRENT_TIMESTAMP)
    ;
END

EVENT_ID=$(psql -qtAXc "SELECT id from event_ WHERE event_.user_id = ${USER_ID} order by id DESC limit 1;")

echo ">> creating some event_items for event id=${EVENT_ID}"
psql -qtAX <<END
INSERT INTO event_item_ 
    (ref_id, event_id, description)
VALUES 
    (decode('$(refidtool generate -o hex -t 3)', 'hex'), ${EVENT_ID}, 'event item 1 description'),
    (decode('$(refidtool generate -o hex -t 3)', 'hex'), ${EVENT_ID}, 'event item 2 description'),
    (decode('$(refidtool generate -o hex -t 3)', 'hex'), ${EVENT_ID}, 'event item 3 description'),
    (decode('$(refidtool generate -o hex -t 3)', 'hex'), ${EVENT_ID}, 'event item 4 description'),
    (decode('$(refidtool generate -o hex -t 3)', 'hex'), ${EVENT_ID}, 'event item 5 description')
    ;
END

EVENT_ITEM_ID=$(psql -qtAXc 'SELECT event_item_.id from event_item_ left join earmark_ on event_item_.id = earmark_.event_item_id where earmark_.event_item_id is NULL limit 1;')
echo ">> creating an earmark for event_item id=${EVENT_ITEM_ID} as user id=${NEXT_USER_ID}"
psql -qtAX <<END
INSERT INTO earmark_ 
    (ref_id, event_item_id, user_id, note)
VALUES 
    (decode('$(refidtool generate -o hex -t 4)', 'hex'), ${EVENT_ITEM_ID}, ${NEXT_USER_ID}, 'i love pickles!')
    ;
END

EVENT_ITEM_ID=$(psql -qtAXc 'SELECT event_item_.id from event_item_ left join earmark_ on event_item_.id = earmark_.event_item_id where earmark_.event_item_id is NULL limit 1;')
echo ">> creating an earmark for event_item id=${EVENT_ITEM_ID} as user id=${NEXT_USER2_ID}"
psql -qtAX <<END
INSERT INTO earmark_ 
    (ref_id, event_item_id, user_id, note)
VALUES 
    (decode('$(refidtool generate -o hex -t 4)', 'hex'), ${EVENT_ITEM_ID}, ${NEXT_USER2_ID}, 'i love pickles!')
    ;
END
