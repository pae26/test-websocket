$(function() {
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
            console.log(stickies);
            let id = stickies[0]['Id'];
            let text = stickies[0]['Text'];
            let x = stickies[0]['Locate_x'];
            let y = stickies[0]['Locate_y'];
            console.log(y);
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
                "Id": id,
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
        console.log("style_str:", style_str);
        let style_list = style_str.match(/[0-9]+/g);
        console.log("style_list:", style_list);
        let location_x = style_list[0];
        let location_y = style_list[1];
        updateSticy(id, location_x, location_y);
    });

    //initSticky();
    loadSticky();
});