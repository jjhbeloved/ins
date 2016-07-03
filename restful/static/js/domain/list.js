/**
 * Created by xiaoxiao on 16/7/2.
 */

var DomainList = {
    init: function () {
        DomainList.ctrl.initEvent();
    },
    ctrl: {
        initEvent: function () {
            $('.btn_delete').on('click', function () {
                var va = $(this).attr("att");
                layer.confirm('are you sure delete the domain?', {icon: 3, title: 'Warning...'}, function (index) {
                    layer.close(index);
                    DomainList.ctrl.deleteAppInfo(va);
                });
            });

        },
        deleteAppInfo: function (id) {
            var dat = {"Id": id};
            layer.load(1);
            $.ajax({
                url: '/jobinfo/delete',
                data: dat,
                dataType: 'json',
                error: function (XMLHttpRequest, textStatus, errorThrown) {
                    layer.closeAll();
                    layer.msg('submit failed.', {
                        icon: 5,
                        time: 2000 //2秒关闭（如果不配置，默认是3秒）
                    }, function () {
                    });

                },
                success: function (data, textStatus, jqXHR) {
                    layer.closeAll();
                    if (data.success == true) {
                        layer.msg(data.message, {
                            icon: 1,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function () {
                            window.location.href = "/domain/list";
                        });

                    } else {
                        layer.msg(data.message, {
                            icon: 5,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function () {
                        });
                    }

                },
                type: 'POST',
                cache: true

            });
        }
    }
};


$(function () {
    DomainList.init();
});

