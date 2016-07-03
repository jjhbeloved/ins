/**
 * Created by xiaoxiao on 16/7/2.
 */

var AppList = {
    init:function(){
        AppList.ctrl.initEvent();
    },
    ctrl:{

        initEvent:function(){

            $('.btn_delete').on('click',function(){
                var va = $(this).attr("att");
                layer.confirm('are you sure delete the app?', {icon: 3, title:'Warning...'}, function(index){
                    layer.close(index);
                    AppList.ctrl.deleteAppInfo(va);
                });
            });

            $("input[type='checkbox']").on('switchChange.bootstrapSwitch', function (event, state) {
                var mySwitch =  $(this);
                var appId = $(this).val();
                mySwitch.bootstrapSwitch('state', !state, true);
               if(state) {
                   layer.confirm('are you sure start the app?', {icon: 3, title:'Warning'}, function(index){
                       layer.close(index);
                       AppList.ctrl.activeApp(mySwitch, state, appId);
                   });
               } else {

                   layer.confirm('are you sure stop the app?', {icon: 3, title:'Warning'}, function(index){
                       layer.close(index);
                       AppList.ctrl.activeApp(mySwitch, state, appId);
                   });
               }

            });
            $("input[type='checkbox']").bootstrapSwitch();
        },

        // active =1 active 0
        activeApp:function(ele, active, appId){

            var dat = {"Id":appId,"active":active};
            layer.load(1);
            $.ajax({
                url:'/app/active',
                data:dat,
                dataType:'json',
                error:function(XMLHttpRequest, textStatus, errorThrown){
                    layer.closeAll();
                    layer.msg('submit failed.', {
                        icon: 5,
                        time: 2000 //2秒关闭（如果不配置，默认是3秒）
                    }, function(){
                    });

                },
                success:function(data, textStatus, jqXHR){
                    layer.closeAll();
                    if(data.success == true) {
                        ele.bootstrapSwitch('state', data.data, true);
                        layer.msg(data.message, {
                            icon: 1,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        });

                    } else {
                        layer.msg(data.message, {
                            icon: 5,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function(){
                        });
                    }

                },
                type:'POST',
                cache:true

            });
        },
        deleteAppInfo:function(id){
            var dat = {"Id":id};
            layer.load(1);
            $.ajax({
                url:'/jobinfo/delete',
                data:dat,
                dataType:'json',
                error:function(XMLHttpRequest, textStatus, errorThrown){
                    layer.closeAll();
                    layer.msg('submit failed.', {
                        icon: 5,
                        time: 2000 //2秒关闭（如果不配置，默认是3秒）
                    }, function(){
                    });

                },
                success:function(data, textStatus, jqXHR){
                    layer.closeAll();
                    if(data.success == true) {
                        layer.msg(data.message, {
                            icon: 1,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function(){
                            window.location.href ="/app/list";
                        });

                    } else {
                        layer.msg(data.message, {
                            icon: 5,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function(){
                        });
                    }

                },
                type:'POST',
                cache:true

            });
        }
    }
};


$(function(){
    AppList.init();
});

