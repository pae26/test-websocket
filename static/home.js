$(function() {
    var socket = null;
    function newSticky(id, text) {
        return '<span class="init-sticky sticky" data-sticky-id="' + id + '" data-color="blue" data-shape="square">' + '<span class="sticky-text">' + text +'</span></span>';
    }

    /*function initSticky() {
        let x = $('.init-sticky').data('x');
        let y = $('.init-sticky').data('y');
        let text = $('.init-sticky').data('text');
    
        let sticky = newSticky(text);
        $('.parent').append(sticky);

        $('.sticky').draggable({
            containment: '.parent',
        });
        $('.sticky').css({
            left: x + 'px',
            top: y + 'px'
        });
    }*/

    function loadSticky() {
        $.ajax({
            dataType: 'json',
            type: 'GET',
            url: '/stickies',
        }).done(function(stickies) {
            console.log('通信成功');
            let id = stickies[0]['Id'];
            let text = stickies[0]['Text'];
            let x = stickies[0]['Locate_x'];
            let y = stickies[0]['Locate_y'];
            let sticky = newSticky(id, text);
            $('.parent').append(sticky);
            $('.sticky').draggable({
                containment: '.slide',
            });
            $('.init-sticky').css({
                left: x + 'px',
                top: y + 'px'
            });
            $('.init-sticky').removeClass('init-sticky');
        }).fail(function(){
            console.log('通信失敗');
        })
    }

    function updateSticy(id, location_x, location_y) {
        $.ajax({
            dataType: 'json',
            contentType: 'application/json',
            type: 'POST',
            url: '/update-sticky',
            data : JSON.stringify({
                "Id": parseInt(id),
                "Locate_x": parseInt(location_x),
                "Locate_y": parseInt(location_y)
            })
        }).done(function() {
            console.log('通信成功');
        }).fail(function() {
            console.log('通信失敗');
        });
    }

    $(document).on('mouseup', '.sticky', function() {
        let id = $(this).data('sticky-id');
        let style_str = $(this).attr('style');
        let style_list = style_str.match(/[0-9]+/g);
        let location_x = style_list[0];
        let location_y = style_list[1];
        updateSticy(id, location_x, location_y);

        if(!socket) {
            alert("エラー: WebSocket通信が行われていません。");
            return false;
        }
        socket.send(id + "," + location_x + "," + location_y);
        return false;
    });

    //initSticky();
    loadSticky();

    if(!window["WebSocket"]) {
        alert("エラー: WebSocketに対応していないブラウザです。");
    } else {
        socket = new WebSocket("ws://localhost:9001/room");
        socket.onopen = function(event) {
            console.log("接続しました:" + event);
	    }
        socket.onclose = function() {
            alert("接続が終了しました。");
        }
        socket.onmessage = function(e) {
            console.log(e.data);
            let socket_slice = e.data.split(',');
            let socket_id = socket_slice[0];
            let socket_x = socket_slice[1];
            let socket_y = socket_slice[2];
            //updateSticy(socket_id, socket_x, socket_y);
            $('[data-sticky-id="'+ socket_id +'"]').animate({
                'left': socket_x + 'px',
                'top': socket_y + 'px'
            })
        }
    }

    let data_label = [];
    for(let i=0;i<47;i++) {
        data_label.push(i.toString());
    }

    let ctx = document.getElementById("chart").getContext('2d');
    let myChart = new Chart(ctx, {
        type: "bar",
        data: {
            labels: data_label,
            datasets: [
                {
                    label: "ページ",
                    data: [10, 22, 10, 9, 12, 6],
                    backgroundColor: "blue",
                },
            ]
        },
        options: {
            responsive: true,
            scales: {
                xAxes: [{
                }],
                yAxes: [{
                    ticks: {
                        min: 0,
                        max: 10,
                        stepSize: 10,
                    },
                }],

            },
            hover: {
                mode: 'single'
            },
        }
    });

    $('.child').on('click', function() {
        let aaa = $(this).parents().find('body');
        console.log('aaa:'+ JSON.stringify(aaa, null, 2));
    });

    let a = true;
    b = 'bbb' + a;
    console.log('b:' + b);
});