myApp = angular.module("reader", []).directive("contenteditable", function() {
	return {
		require: 'ngModel',
		link: function(scope, elm, attrs, ctrl) {
			// view -> model
			elm.bind('blur', function() {
				scope.$apply(function() {
					ctrl.$setViewValue(elm.html());
				});
			});
			// model -> view
			ctrl.$render = function() {
				elm.html(ctrl.$viewValue);
			};

			// load init value from DOM
			ctrl.$setViewValue(elm.html());
		}
	}
})

function FeedAdminCtrl($scope, $http) {
	$scope.showError = false
	$scope.showAdd = false

	function log(msg) {
		if (console.log) {
			console.log(msg)
		}
	}

	$scope.refresh = function() {
		log("Refreshing")
		$scope.showError = true
		$scope.error = "Fetching"

		$http.get('/feed').success(function (data) {
			log("Refresh Status " + data.Status)
			if (data.Status != "ok") {
				$scope.showError = true
				$scope.error = data.Message
			} else {
				$scope.showError = false
				$scope.feeds = data.Feeds
			}
		})	
	}

	$scope.addId = "id"
	$scope.addUrl = "url"
	$scope.toggleAdd = function() {
		if ($scope.showAdd) {
			$scope.showAdd = false
		} else {
			$scope.showAdd = true
		}
	}
	$scope.addFeed = function() {
		log("Adding")
		var data = JSON.stringify({ Id: $scope.addId, Url: $scope.addUrl })
		var promise = $http.post("/feed/", data)
		promise.success(function (data) {
			log("Add Status " + data.Status)
			if (data.Status != "ok") {
				$scope.showError = true
				$scope.error = data.Message
			} else {
				$scope.showError = false
				$scope.showAdd = false
				$scope.addId = "id"
				$scope.addUrl = "url"
				$scope.refresh()
			}
		})
	}
	$scope.refreshFeed = function(feed) {
		log("Refreshing " + feed.Id)
		$http.post("/feed/refresh/" + feed.Id).success(function (data) {
			log("Refresh Feed Status " + data.Status)
			if (data.Status != "ok") {
				$scope.showError = true
				$scope.error = data.Message
			} else {
				$scope.showError = false
				if (console.log) {
					console.log(data.Message)
				}
			}
		})
	}

	$scope.refresh()
}
FeedAdminCtrl.$inject = ['$scope', '$http']
