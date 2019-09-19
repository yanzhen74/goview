var __extends=this&&this.__extends||function(){var a=function(c,d){return a=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(a,c){a.__proto__=c}||function(a,c){for(var b in c)c.hasOwnProperty(b)&&(a[b]=c[b])},a(c,d)};return function(c,d){function b(){this.constructor=c}a(c,d),c.prototype=null===d?Object.create(d):(b.prototype=d.prototype,new b)}}(),__awaiter=this&&this.__awaiter||function(a,b,c,d){function e(a){return a instanceof c?a:new c(function(b){b(a)})}return new(c||(c=Promise))(function(c,f){function g(a){try{i(d.next(a))}catch(a){f(a)}}function h(a){try{i(d["throw"](a))}catch(a){f(a)}}function i(a){a.done?c(a.value):e(a.value).then(g,h)}i((d=d.apply(a,b||[])).next())})},__generator=this&&this.__generator||function(a,b){function c(a){return function(b){return d([a,b])}}function d(c){if(e)throw new TypeError("Generator is already executing.");for(;k;)try{if(e=1,h&&(i=2&c[0]?h["return"]:c[0]?h["throw"]||((i=h["return"])&&i.call(h),0):h.next)&&!(i=i.call(h,c[1])).done)return i;switch((h=0,i)&&(c=[2&c[0],i.value]),c[0]){case 0:case 1:i=c;break;case 4:return k.label++,{value:c[1],done:!1};case 5:k.label++,h=c[1],c=[0];continue;case 7:c=k.ops.pop(),k.trys.pop();continue;default:if((i=k.trys,!(i=0<i.length&&i[i.length-1]))&&(6===c[0]||2===c[0])){k=0;continue}if(3===c[0]&&(!i||c[1]>i[0]&&c[1]<i[3])){k.label=c[1];break}if(6===c[0]&&k.label<i[1]){k.label=i[1],i=c;break}if(i&&k.label<i[2]){k.label=i[2],k.ops.push(c);break}i[2]&&k.ops.pop(),k.trys.pop();continue;}c=b.call(a,k)}catch(a){c=[6,a],h=0}finally{e=i=0}if(5&c[0])throw c[1];return{value:c[0]?c[1]:void 0,done:!0}}var e,h,i,j,k={label:0,sent:function(){if(1&i[0])throw i[1];return i[1]},trys:[],ops:[]};return j={next:c(0),throw:c(1),return:c(2)},"function"==typeof Symbol&&(j[Symbol.iterator]=function(){return this}),j},isBrowser="undefined"!=typeof window,_fetch="undefined"==typeof fetch?void 0:fetch;isBrowser?WebSocket=window.WebSocket:(WebSocket=require("ws"),_fetch=require("node-fetch"));var OnNamespaceConnect="_OnNamespaceConnect",OnNamespaceConnected="_OnNamespaceConnected",OnNamespaceDisconnect="_OnNamespaceDisconnect",OnRoomJoin="_OnRoomJoin",OnRoomJoined="_OnRoomJoined",OnRoomLeave="_OnRoomLeave",OnRoomLeft="_OnRoomLeft",OnAnyEvent="_OnAnyEvent",OnNativeMessage="_OnNativeMessage",ackBinary="M",ackIDBinary="A",ackNotOKBinary="H",waitIsConfirmationPrefix="#",waitComesFromClientPrefix="$";function isSystemEvent(a){return!(a!==OnNamespaceConnect&&a!==OnNamespaceConnected&&a!==OnNamespaceDisconnect&&a!==OnRoomJoin&&a!==OnRoomJoined&&a!==OnRoomLeave&&a!==OnRoomLeft)}function isEmpty(a){return!(void 0!==a)||!(null!==a)||(""==a||"string"==typeof a||a instanceof String?0===a.length||""===a:!!(a instanceof Error)&&isEmpty(a.message))}var Message=function(){function a(){}return a.prototype.isConnect=function(){return this.Event==OnNamespaceConnect||!1},a.prototype.isDisconnect=function(){return this.Event==OnNamespaceDisconnect||!1},a.prototype.isRoomJoin=function(){return this.Event==OnRoomJoin||!1},a.prototype.isRoomLeft=function(){return this.Event==OnRoomLeft||!1},a.prototype.isWait=function(){return!isEmpty(this.wait)&&(this.wait[0]==waitIsConfirmationPrefix||this.wait[0]==waitComesFromClientPrefix||!1)},a.prototype.unmarshal=function(){return JSON.parse(this.Body)},a}();function marshal(a){return JSON.stringify(a)}var messageSeparator=";",messageFieldSeparatorReplacement="@%!semicolon@%!",validMessageSepCount=7,trueString="1",falseString="0",escapeRegExp=/;/g;function escapeMessageField(a){return isEmpty(a)?"":a.replace(escapeRegExp,messageFieldSeparatorReplacement)}var unescapeRegExp=new RegExp(messageFieldSeparatorReplacement,"g");function unescapeMessageField(a){return isEmpty(a)?"":a.replace(unescapeRegExp,messageSeparator)}var replyError=function(a){function b(c){var d=a.call(this,c)||this;return d.name="replyError",Error.captureStackTrace(d,b),Object.setPrototypeOf(d,b.prototype),d}return __extends(b,a),b}(Error);function reply(a){return new replyError(a)}function isReply(a){return a instanceof replyError}function serializeMessage(a){if(a.IsNative&&isEmpty(a.wait))return a.Body;var b=falseString,c=falseString,d=a.Body||"";return isEmpty(a.Err)||(d=a.Err.message,!isReply(a.Err)&&(b=trueString)),a.isNoOp&&(c=trueString),[a.wait||"",escapeMessageField(a.Namespace),escapeMessageField(a.Room),escapeMessageField(a.Event),b,c,d].join(messageSeparator)}function splitN(a,b,c){if(0==c)return[a];var d=a.split(b,c);if(d.length==c){var e=d.join(b)+b;return d.push(a.substr(e.length)),d}return[a]}var textDecoder=new TextDecoder("utf-8"),messageSeparatorCharCode=";".charCodeAt(0);function deserializeMessage(a,b){var c=new Message;if(0==a.length)return c.isInvalid=!0,c;var d,e=a instanceof ArrayBuffer;if(e){for(var f=new Uint8Array(a),g=1,h=0,j=0;j<f.length;j++)f[j]==messageSeparatorCharCode&&(g++,h=j);if(g!=validMessageSepCount)return c.isInvalid=!0,c;d=splitN(textDecoder.decode(f.slice(0,h)),messageSeparator,validMessageSepCount-2),d.push(a.slice(h+1,a.length)),c.SetBinary=!0}else d=splitN(a,messageSeparator,validMessageSepCount-1);if(d.length!=validMessageSepCount)return b?(c.Event=OnNativeMessage,c.Body=a):c.isInvalid=!0,c;c.wait=d[0],c.Namespace=unescapeMessageField(d[1]),c.Room=unescapeMessageField(d[2]),c.Event=unescapeMessageField(d[3]),c.isError=d[4]==trueString||!1,c.isNoOp=d[5]==trueString||!1;var k=d[6];return isEmpty(k)?c.Body="":c.isError?c.Err=new Error(k):c.Body=k,c.isInvalid=!1,c.IsForced=!1,c.IsLocal=!1,c.IsNative=b&&c.Event==OnNativeMessage||!1,c}function genWait(){if(!isBrowser){var a=process.hrtime();return waitComesFromClientPrefix+1e9*a[0]+a[1]}var b=window.performance.now();return waitComesFromClientPrefix+b.toString()}function genWaitConfirmation(a){return waitIsConfirmationPrefix+a}function genEmptyReplyToWait(a){return a+messageSeparator.repeat(validMessageSepCount-1)}var Room=function(){function a(a,b){this.nsConn=a,this.name=b}return a.prototype.emit=function(a,b){var c=new Message;return c.Namespace=this.nsConn.namespace,c.Room=this.name,c.Event=a,c.Body=b,this.nsConn.conn.write(c)},a.prototype.leave=function(){var a=new Message;return a.Namespace=this.nsConn.namespace,a.Room=this.name,a.Event=OnRoomLeave,this.nsConn.askRoomLeave(a)},a}(),NSConn=function(){function a(a,b,c){this.conn=a,this.namespace=b,this.events=c,this.rooms=new Map}return a.prototype.emit=function(a,b){var c=new Message;return c.Namespace=this.namespace,c.Event=a,c.Body=b,this.conn.write(c)},a.prototype.ask=function(a,b){var c=new Message;return c.Namespace=this.namespace,c.Event=a,c.Body=b,this.conn.ask(c)},a.prototype.joinRoom=function(a){return __awaiter(this,void 0,void 0,function(){return __generator(this,function(b){switch(b.label){case 0:return[4,this.askRoomJoin(a)];case 1:return[2,b.sent()];}})})},a.prototype.room=function(a){return this.rooms.get(a)},a.prototype.leaveAll=function(){return __awaiter(this,void 0,void 0,function(){var a,b=this;return __generator(this,function(){return a=new Message,a.Namespace=this.namespace,a.Event=OnRoomLeft,a.IsLocal=!0,this.rooms.forEach(function(c,d){return __awaiter(b,void 0,void 0,function(){var b;return __generator(this,function(c){switch(c.label){case 0:a.Room=d,c.label=1;case 1:return c.trys.push([1,3,,4]),[4,this.askRoomLeave(a)];case 2:return c.sent(),[3,4];case 3:return b=c.sent(),[2,b];case 4:return[2];}})})}),[2,null]})})},a.prototype.forceLeaveAll=function(a){var b=this,c=new Message;c.Namespace=this.namespace,c.Event=OnRoomLeave,c.IsForced=!0,c.IsLocal=a,this.rooms.forEach(function(a,d){c.Room=d,fireEvent(b,c),b.rooms.delete(d),c.Event=OnRoomLeft,fireEvent(b,c),c.Event=OnRoomLeave})},a.prototype.disconnect=function(){var a=new Message;return a.Namespace=this.namespace,a.Event=OnNamespaceDisconnect,this.conn.askDisconnect(a)},a.prototype.askRoomJoin=function(a){var b=this;return new Promise(function(c,d){return __awaiter(b,void 0,void 0,function(){var b,e,f,g;return __generator(this,function(h){switch(h.label){case 0:if(b=this.rooms.get(a),void 0!==b)return c(b),[2];e=new Message,e.Namespace=this.namespace,e.Room=a,e.Event=OnRoomJoin,e.IsLocal=!0,h.label=1;case 1:return h.trys.push([1,3,,4]),[4,this.conn.ask(e)];case 2:return h.sent(),[3,4];case 3:return f=h.sent(),d(f),[2];case 4:return(g=fireEvent(this,e),!isEmpty(g))?(d(g),[2]):(b=new Room(this,a),this.rooms.set(a,b),e.Event=OnRoomJoined,fireEvent(this,e),c(b),[2]);}})})})},a.prototype.askRoomLeave=function(a){return __awaiter(this,void 0,void 0,function(){var b,c;return __generator(this,function(d){switch(d.label){case 0:if(!this.rooms.has(a.Room))return[2,ErrBadRoom];d.label=1;case 1:return d.trys.push([1,3,,4]),[4,this.conn.ask(a)];case 2:return d.sent(),[3,4];case 3:return b=d.sent(),[2,b];case 4:return(c=fireEvent(this,a),!isEmpty(c))?[2,c]:(this.rooms.delete(a.Room),a.Event=OnRoomLeft,fireEvent(this,a),[2,null]);}})})},a.prototype.replyRoomJoin=function(a){if(!(isEmpty(a.wait)||a.isNoOp)){if(!this.rooms.has(a.Room)){var b=fireEvent(this,a);if(!isEmpty(b))return a.Err=b,void this.conn.write(a);this.rooms.set(a.Room,new Room(this,a.Room)),a.Event=OnRoomJoined,fireEvent(this,a)}this.conn.writeEmptyReply(a.wait)}},a.prototype.replyRoomLeave=function(a){return isEmpty(a.wait)||a.isNoOp?void 0:this.rooms.has(a.Room)?void(fireEvent(this,a),this.rooms.delete(a.Room),this.conn.writeEmptyReply(a.wait),a.Event=OnRoomLeft,fireEvent(this,a)):void this.conn.writeEmptyReply(a.wait)},a}();function fireEvent(a,b){return a.events.has(b.Event)?a.events.get(b.Event)(a,b):a.events.has(OnAnyEvent)?a.events.get(OnAnyEvent)(a,b):null}function isNull(a){return null===a||a===void 0||"undefined"==typeof a}function resolveNamespaces(a,b){if(isNull(a))return isNull(b)||b("connHandler is empty."),null;var c=new Map,d=new Map,e=0;if(Object.keys(a).forEach(function(b){e++;var f=a[b];if(f instanceof Function)d.set(b,f);else if(f instanceof Map)c.set(b,f);else{var g=new Map;Object.keys(f).forEach(function(a){g.set(a,f[a])}),c.set(b,g)}}),0<d.size){if(e!=d.size)return isNull(b)||b("all keys of connHandler should be events, mix of namespaces and event callbacks is not supported "+d.size+" vs total "+e),null;c.set("",d)}return c}function getEvents(a,b){return a.has(b)?a.get(b):null}var URLParamAsHeaderPrefix="X-Websocket-Header-";function parseHeadersAsURLParameters(a,b){if(isNull(a))return b;for(var c in a)if(a.hasOwnProperty(c)){var d=a[c];c=encodeURIComponent(URLParamAsHeaderPrefix+c),d=encodeURIComponent(d);var e=c+"="+d;b=-1==b.indexOf("?")?-1==b.indexOf("#")?b+"?"+e:b.split("#")[0]+"?"+e+"#"+b.split("#")[1]:b.split("?")[0]+"?"+e+"&"+b.split("?")[1]}return b}function dial(a,b,c){return _dial(a,b,0,c)}var websocketReconnectHeaderKey="X-Websocket-Reconnect";function _dial(a,b,c,d){if(isBrowser&&0==a.indexOf("/")){var e="https:"==document.location.protocol?"wss":"ws",f=document.location.port?":"+document.location.port:"";a=e+"://"+document.location.hostname+f+a}return-1==a.indexOf("ws")&&(a="ws://"+a),new Promise(function(e,f){WebSocket||f("WebSocket is not accessible through this browser.");var g=resolveNamespaces(b,f);if(!isNull(g)){isNull(d)&&(d={}),isNull(d.headers)&&(d.headers={});var h=d.reconnect?d.reconnect:0;0<c&&0<h?d.headers[websocketReconnectHeaderKey]=c.toString():!isNull(d.headers[websocketReconnectHeaderKey])&&delete d.headers[websocketReconnectHeaderKey];var i=makeWebsocketConnection(a,d),j=new Conn(i,g);j.reconnectTries=c,i.binaryType="arraybuffer",i.onmessage=function(a){var b=j.handle(a);return isEmpty(b)?void(j.isAcknowledged()&&e(j)):void f(b)},i.onopen=function(){i.send(ackBinary)},i.onerror=function(a){j.close(),f(a)},i.onclose=function(){if(j.isClosed());else{if(i.onmessage=void 0,i.onopen=void 0,i.onerror=void 0,i.onclose=void 0,0>=h)return j.close(),null;var c=new Map;j.connectedNamespaces.forEach(function(a,b){var d=[];!isNull(a.rooms)&&0<a.rooms.size&&a.rooms.forEach(function(a,b){d.push(b)}),c.set(b,d)}),j.close(),whenResourceOnline(a,h,function(g){_dial(a,b,g,d).then(function(a){return isNull(e)||"function () { [native code] }"==e.toString()?void c.forEach(function(b,c){a.connect(c).then(function(a){return function(b){a.forEach(function(a){b.joinRoom(a)})}}(b))}):void e(a)}).catch(f)})}return null}}})}function makeWebsocketConnection(a,b){return isBrowser&&!isNull(b)?(b.headers&&(a=parseHeadersAsURLParameters(b.headers,a)),b.protocols?new WebSocket(a,b.protocols):new WebSocket(a)):new WebSocket(a,b)}function whenResourceOnline(a,b,c){var d=a.replace(/(ws)(s)?\:\/\//,"http$2://"),e=1,f={method:"HEAD",mode:"no-cors"},g=function(){_fetch(d,f).then(function(){c(e)}).catch(function(){e++,setTimeout(function(){g()},b)})};setTimeout(g,b)}var ErrInvalidPayload=new Error("invalid payload"),ErrBadNamespace=new Error("bad namespace"),ErrBadRoom=new Error("bad room"),ErrClosed=new Error("use of closed connection"),ErrWrite=new Error("write closed");function isCloseError(a){return!(!a||isEmpty(a.message))&&0<=a.message.indexOf("[-1] write closed")}var Conn=function(){function a(a,b){this.conn=a,this.reconnectTries=0,this._isAcknowledged=!1,this.namespaces=b;var c=b.has("");this.allowNativeMessages=c&&b.get("").has(OnNativeMessage),this.queue=[],this.waitingMessages=new Map,this.connectedNamespaces=new Map,this.closed=!1}return a.prototype.wasReconnected=function(){return 0<this.reconnectTries},a.prototype.isAcknowledged=function(){return this._isAcknowledged},a.prototype.handle=function(a){if(!this._isAcknowledged){var b=this.handleAck(a.data);return null==b?(this._isAcknowledged=!0,this.handleQueue()):this.conn.close(),b}return this.handleMessage(a.data)},a.prototype.handleAck=function(a){var b=a[0];switch(b){case ackIDBinary:var c=a.slice(1);this.ID=c;break;case ackNotOKBinary:var d=a.slice(1);return new Error(d);default:return this.queue.push(a),null;}},a.prototype.handleQueue=function(){var a=this;null==this.queue||0==this.queue.length||this.queue.forEach(function(b,c){a.queue.splice(c,1),a.handleMessage(b)})},a.prototype.handleMessage=function(a){var b=deserializeMessage(a,this.allowNativeMessages);if(b.isInvalid)return ErrInvalidPayload;if(b.IsNative&&this.allowNativeMessages){var c=this.namespace("");return fireEvent(c,b)}if(b.isWait()){var d=this.waitingMessages.get(b.wait);if(null!=d)return void d(b)}var e=this.namespace(b.Namespace);switch(b.Event){case OnNamespaceConnect:this.replyConnect(b);break;case OnNamespaceDisconnect:this.replyDisconnect(b);break;case OnRoomJoin:if(void 0!==e){e.replyRoomJoin(b);break}case OnRoomLeave:if(void 0!==e){e.replyRoomLeave(b);break}default:if(void 0===e)return ErrBadNamespace;b.IsLocal=!1;var f=fireEvent(e,b);if(!isEmpty(f))return b.Err=f,this.write(b),f;}return null},a.prototype.connect=function(a){return this.askConnect(a)},a.prototype.waitServerConnect=function(a){var b=this;return isNull(this.waitServerConnectNotifiers)&&(this.waitServerConnectNotifiers=new Map),new Promise(function(c){return __awaiter(b,void 0,void 0,function(){var b=this;return __generator(this,function(){return this.waitServerConnectNotifiers.set(a,function(){b.waitServerConnectNotifiers.delete(a),c(b.namespace(a))}),[2]})})})},a.prototype.namespace=function(a){return this.connectedNamespaces.get(a)},a.prototype.replyConnect=function(a){if(!(isEmpty(a.wait)||a.isNoOp)){var b=this.namespace(a.Namespace);if(void 0!==b)return void this.writeEmptyReply(a.wait);var c=getEvents(this.namespaces,a.Namespace);return isNull(c)?(a.Err=ErrBadNamespace,void this.write(a)):void(b=new NSConn(this,a.Namespace,c),this.connectedNamespaces.set(a.Namespace,b),this.writeEmptyReply(a.wait),a.Event=OnNamespaceConnected,fireEvent(b,a),!isNull(this.waitServerConnectNotifiers)&&0<this.waitServerConnectNotifiers.size&&this.waitServerConnectNotifiers.has(a.Namespace)&&this.waitServerConnectNotifiers.get(a.Namespace)())}},a.prototype.replyDisconnect=function(a){if(!(isEmpty(a.wait)||a.isNoOp)){var b=this.namespace(a.Namespace);return void 0===b?void this.writeEmptyReply(a.wait):void(b.forceLeaveAll(!0),this.connectedNamespaces.delete(a.Namespace),this.writeEmptyReply(a.wait),fireEvent(b,a))}},a.prototype.ask=function(a){var b=this;return new Promise(function(c,d){return b.isClosed()?void d(ErrClosed):(a.wait=genWait(),b.waitingMessages.set(a.wait,function(a){return a.isError?void d(a.Err):void c(a)}),!b.write(a))?void d(ErrWrite):void 0})},a.prototype.askConnect=function(a){var b=this;return new Promise(function(c,d){return __awaiter(b,void 0,void 0,function(){var b,e,f,g,h;return __generator(this,function(i){switch(i.label){case 0:if(b=this.namespace(a),void 0!==b)return c(b),[2];if(e=getEvents(this.namespaces,a),isNull(e))return d(ErrBadNamespace),[2];if(f=new Message,f.Namespace=a,f.Event=OnNamespaceConnect,f.IsLocal=!0,b=new NSConn(this,a,e),g=fireEvent(b,f),!isEmpty(g))return d(g),[2];i.label=1;case 1:return i.trys.push([1,3,,4]),[4,this.ask(f)];case 2:return i.sent(),[3,4];case 3:return h=i.sent(),d(h),[2];case 4:return this.connectedNamespaces.set(a,b),f.Event=OnNamespaceConnected,fireEvent(b,f),c(b),[2];}})})})},a.prototype.askDisconnect=function(a){return __awaiter(this,void 0,void 0,function(){var b,c;return __generator(this,function(d){switch(d.label){case 0:if(b=this.namespace(a.Namespace),void 0===b)return[2,ErrBadNamespace];d.label=1;case 1:return d.trys.push([1,3,,4]),[4,this.ask(a)];case 2:return d.sent(),[3,4];case 3:return c=d.sent(),[2,c];case 4:return b.forceLeaveAll(!0),this.connectedNamespaces.delete(a.Namespace),a.IsLocal=!0,[2,fireEvent(b,a)];}})})},a.prototype.isClosed=function(){return this.closed},a.prototype.write=function(a){if(this.isClosed())return!1;if(!a.isConnect()&&!a.isDisconnect()){var b=this.namespace(a.Namespace);if(void 0===b)return!1;if(!isEmpty(a.Room)&&!a.isRoomJoin()&&!a.isRoomLeft()&&!b.rooms.has(a.Room))return!1}return this.conn.send(serializeMessage(a)),!0},a.prototype.writeEmptyReply=function(a){this.conn.send(genEmptyReplyToWait(a))},a.prototype.close=function(){var a=this;if(!this.closed){var b=new Message;b.Event=OnNamespaceDisconnect,b.IsForced=!0,b.IsLocal=!0,this.connectedNamespaces.forEach(function(c){c.forceLeaveAll(!0),b.Namespace=c.namespace,fireEvent(c,b),a.connectedNamespaces.delete(c.namespace)}),this.waitingMessages.clear(),this.closed=!0,this.conn.readyState===this.conn.OPEN&&this.conn.close()}},a}();(function(){var a={dial:dial,isSystemEvent:isSystemEvent,OnNamespaceConnect:OnNamespaceConnect,OnNamespaceConnected:OnNamespaceConnected,OnNamespaceDisconnect:OnNamespaceDisconnect,OnRoomJoin:OnRoomJoin,OnRoomJoined:OnRoomJoined,OnRoomLeave:OnRoomLeave,OnRoomLeft:OnRoomLeft,OnAnyEvent:OnAnyEvent,OnNativeMessage:OnNativeMessage,Message:Message,Room:Room,NSConn:NSConn,Conn:Conn,ErrInvalidPayload:ErrInvalidPayload,ErrBadNamespace:ErrBadNamespace,ErrBadRoom:ErrBadRoom,ErrClosed:ErrClosed,ErrWrite:ErrWrite,isCloseError:isCloseError,reply:reply,marshal:marshal};if("undefined"!=typeof exports)exports=a,module.exports=a;else{var b="object"==typeof self&&self.self===self&&self||"object"==typeof global&&global.global===global&&global;b.neffos=a}})();