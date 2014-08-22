api_endpoint = 'https://api.hipchat.com'
api_key = "<your hipchat api token>"

// https://www.hipchat.com/docs/api
var postMessage = function(params, callback) {
  var data = {
    room_id: params.room_id || params.room,
    from: (params.from || '--'),
    message: params.message,
    notify: 1,
    color: (params.color || 'yellow'),
    message_format: 'text',
  };
  return $.ajax({
    type: 'post',
    url: api_endpoint + '/v1/rooms/message?auth_token=' + api_key,
    data: data,
    success: function(body, status, res) { callback(body, status, res) },
    failure: function() { console.log(arguments) }
  });
};

function main(params, ctx, done) {
  var params = {
      room: params.room_id,
      message: params.message,
  }
  postMessage(params, function(data) {
    done({params:params, ctx:ctx, data:data});
  })
}
