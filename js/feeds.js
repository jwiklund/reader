function FeedCtrl($scope, $html) {
	$html.get('/feed').success(function (data) {
		if (data.Status != "Ok") {
			$scope.error = data.Message
		} else {
			$scope.feeds = data.Data
		}
	})
}
FeedCtrl.$inject = ['$scope', '$http']
