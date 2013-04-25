function FeedCtrl($scope, $html) {
	$scope.showError = false

	function log(msg) {
		if (console.log) {
			console.log(msg)
		}
	}

	$scope.refresh = function() {
		log("Refreshing")
		$scope.showError = true
		$scope.error = "Fetching"
	}

	$scope.refresh()
}
FeedCtrl.$inject = ['$scope', '$http']
