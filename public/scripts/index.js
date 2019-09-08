  layui.use(['element', 'layer', 'jquery', 'tree', 'util'], function(){
    
    var $ = layui.$
    ,tree = layui.tree
    ,layer = layui.layer
    ,util = layui.util
    ,index = 100;
    
    //数据源
    var data1 = [{
      title: 'gcyc'
      ,id: 1
      ,spread: true
      ,children: [{
        title: 'gcyctd'
        ,id: 3
        ,spread: true
        //,href: 'https://www.layui.com/doc/'
        ,href: "javascript:;" 
        ,class: "site_demo_active"
        ,data_type: "tabAdd"
        ,data_title: "pk25"
        ,data_id: "1"
        ,data_url: "https://www.layui.com/doc/"
        }],
      title: 'sgyc'
      ,id: 2
      ,spread: true
      ,children: [{
        title: 'pk25'
        ,id: 5
        ,spread: true
        ,href: "javascript:;" 
        ,class: "site_demo_active"
        ,data_type: "tabAdd"
        ,data_title: "pk25"
        ,data_id: "2"
        ,data_url: "https://www.layui.com/doc/"
        ,children: [{
          title: '三级2_1_1'
          ,id: 11
            ,href: "javascript:;" 
            ,class: "site_demo_active"
            ,data_type: "tabAdd"
            ,data_title: "pk25"
            ,data_id: "3"
            ,data_url: "https://www.layui.com/doc/"
        },{
          title: '三级2_1_2'
          ,id: 12
        }]
      },{
        title: '二级2_2'
        ,id: 6
        ,checked: true
        ,children: [{
          title: '三级2_2_1'
          ,id: 13
        },{
          title: '三级2_2_2'
          ,id: 14
          ,disabled: true
        }]
      }]
    }];

    
    tree.render({
      elem: '#test2'
      ,data: data1
      //,expandClick: false
      ,showLine: true 
      ,click: function(obj, state){
        console.log(obj);
        var dataid = $(obj)[0]["data"];
        console.log(dataid["id"]);
        if (dataid["class"] == "site_demo_active") {
            console.log(dataid["data_url"]);
            //这时会判断右侧.layui-tab-title属性下的有lay-id属性的li的数目，即已经打开的tab项数目
            if ($(".layui-tab-title li[lay-id]").length <= 0) {
                //如果比零小，则直接打开新的tab项
                active.tabAdd(dataid["data_url"], dataid["data_id"], dataid["data_title"]);
            } else {
                //否则判断该tab项是否以及存在
                var isData = false; //初始化一个标志，为false说明未打开该tab项 为true则说明已有
                $.each($(".layui-tab-title li[lay-id]"), function () {
                    //如果点击左侧菜单栏所传入的id 在右侧tab项中的lay-id属性可以找到，则说明该tab项已经打开
                    if ($(this).attr("lay-id") == dataid["data_id"]) {
                        isData = true;
                    }
                })
                if (isData == false) {
                    //标志为false 新增一个tab项
                    active.tabAdd(dataid["data_url"], dataid["data_id"], dataid["data_title"]);
                }
            }
            //最后不管是否新增tab，最后都转到要打开的选项页面上
            active.tabChange(dataid["data_id"]);
        }
      }
      ,oncheck: function(obj, checked, child){
        if(checked){
          console.log(obj[0]);
        }
      }
      ,onsearch: function(data, num){
        console.log(num);
      }
      ,dragstart: function(obj, parent){
        console.log(obj, parent);
      }
      ,dragend: function(state, obj, target){
        console.log(state, obj, target);
      }
    });
     
    var element = layui.element;
    // var layer = layui.layer;
    var $ = layui.$;
    // 配置tab实践在下面无法获取到菜单元素
    $('.site_demo_active').on('click', function () {
        var dataid = $(this);
        //这时会判断右侧.layui-tab-title属性下的有lay-id属性的li的数目，即已经打开的tab项数目
        if ($(".layui-tab-title li[lay-id]").length <= 0) {
            //如果比零小，则直接打开新的tab项
            active.tabAdd(dataid["data_url"], dataid["data_id"], dataid["data_title"]);
        } else {
            //否则判断该tab项是否以及存在
            var isData = false; //初始化一个标志，为false说明未打开该tab项 为true则说明已有
            $.each($(".layui-tab-title li[lay-id]"), function () {
                //如果点击左侧菜单栏所传入的id 在右侧tab项中的lay-id属性可以找到，则说明该tab项已经打开
                if ($(this).attr("lay-id") == dataid["data_id"]) {
                    isData = true;
                }
            })
            if (isData == false) {
                //标志为false 新增一个tab项
                active.tabAdd(dataid["data_url"], dataid["data_id"], dataid["data_title"]);
            }
        }
        //最后不管是否新增tab，最后都转到要打开的选项页面上
        active.tabChange(dataid["data_id"]);
    });

    var active = {
        //在这里给active绑定几项事件，后面可通过active调用这些事件
        tabAdd: function (url, id, name) {
            //新增一个Tab项 传入三个参数，分别对应其标题，tab页面的地址，还有一个规定的id，是标签中data_id的属性值
            //关于tabAdd的方法所传入的参数可看layui的开发文档中基础方法部分
            element.tabAdd('demo', {
                title: name,
                content: '<iframe data_frameid="' + id + '" scrolling="auto" frameborder="0" src="' + url + '" style="width:100%;height:99%;"></iframe>',
                id: id //规定好的id
            })
            FrameWH();  //计算ifram层的大小
        },
        tabChange: function (id) {
            //切换到指定Tab项
            element.tabChange('demo', id); //根据传入的id传入到指定的tab项
        },
        tabDelete: function (id) {
            element.tabDelete("demo", id);//删除
        }
    };
    function FrameWH() {
        var h = $(window).height();
        $("iframe").css("height",h+"px");
    }
});