//
// server code example
//
function main(params, ctx) {
    console.log("hello world")
    return {
      "versions": {
        "storage": Kii.getSDKVersion(),
        "analytics": KiiAnalytics.getSDKVersion(),
      },
      "now": (new Date)
    }
}
