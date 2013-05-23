function FeedCtrl($scope, $http) {
	$scope.showError = false

	function log(msg) {
		if (console.log) {
			console.log(msg)
		}
	}

	function formatItems(items) {
		result = []
		first = true
		for (index in items) {
			var item = items[index]
			var resultItem = { markClass: "newItem", alternative: false }
			if (first) {
				first = false
				resultItem.markClass = "curItem"
			}
			if ('Description' in item) {
				resultItem.content = item.Description
			} else if ('Title' in item) {
				resultItem.alternative = true
				resultItem.text = item.Title
				resultItem.url = item.Url
			} else {
				resultItem.alternative = true
				resultItem.text = item.Id
				resultItem.url = item.Url
			}
			result.push(resultItem)
		}
		return result
	}

	$scope.refresh = function() {
		log("Refreshing")
		$scope.showError = true
		$scope.error = "Fetching"

		$http.get('/feed/user/jwiklund/all').success(function (data) {
			log("Refresh Status " + data.Status)
			if (data.Status != "ok") {
				$scope.showError = true
				$scope.error = data.Message
			} else {
				$scope.showError = false
				$scope.items = formatItems(data.Items)
			}
		})
		$scope.showError = true
		$scope.error = "Fetching"
	}

	$scope.handleSpace = function() {
		console.log("SPACE")
	}

	$scope.refresh()
}
FeedCtrl.$inject = ['$scope', '$http']
