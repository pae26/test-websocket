$(function() {
    /*$('.child').draggable({
        containment: '.parent',
    });*/

    function newSticky(text) {
        return '<span class="sticky" data-color="blue" data-shape="square">' + '<span class="sticky-text">' + text +'</span></span>';
    }

    function initSticky() {
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
    }

    initSticky();
});