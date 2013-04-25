function FeedCtrl($scope, $html) {
	$scope.showError = false

	function log(msg) {
		if (console.log) {
			console.log(msg)
		}
	}

	function formatItems(items) {
		result = []
		for (item in items) {
			if ('Description' in items[item]) {
				result.push({ alternative: false, content: items[item].Description })
			} else if ('Title' in items[item]) {
				result.push({ alternative: true, text: items[item].Title, url: items[item].Url})
			} else {
				result.push({ alternative: true, text: items[item].Id, url: items[item].Url})
			}
		}
		return result
	}

	$scope.refresh = function() {
		log("Refreshing")
		$scope.showError = true
		$scope.error = "Fetching"

		$html.get('/read').success(function (data) {
			log("Refresh Status " + data.Status)
			if (data.Status != "Ok") {
				$scope.showError = true
				$scope.error = data.Status
			} else {
				$scope.showError = false
				$scope.items = formatItems(data.Data)
			}
		})
		$scope.showError = true
		$scope.error = "Fetching"
	}

	$scope.refresh()
}
FeedCtrl.$inject = ['$scope', '$http']
