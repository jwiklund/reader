var module = angular.module('reader', []);

/**
 * Directive for binding keyboard shortcuts.
 *
 * When a key press matches one of the key bindings, the associated expression is executed.
 */
module.directive('wKeydown', function() {
  return function(scope, elm, attr) {
    elm.bind('keydown', function(e) {
      switch (e.keyCode) {
        case 32: // Space
          e.preventDefault();
          return scope.$apply(attr.wSpace);
      }
    });
  };
});

/**
 * Service that is in charge of scrolling in the app.
 */
module.factory('scroll', function($timeout) {
	return {
		toNext: function() {
			var cur = $('.item.curItem')
			$('body').scrollTop(cur.offset().top + cur.height() - 60)
		},
		toTop: function() {
			$('body').scrollTop(0)
		}
	}
})

function FeedCtrl($scope, $http, scroll) {
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
				$scope.current = 0
				scroll.toTop()
			}
		})
		$scope.showError = true
		$scope.error = "Fetching"
	}

	$scope.handleSpace = function() {
		if ($scope.items && $scope.current < $scope.items.length - 1) {
			scroll.toNext()
			$scope.items[$scope.current].markClass = ""
			$scope.items[$scope.current + 1].markClass = "curItem"
			$scope.current = $scope.current + 1
		}
	}

	$scope.refresh()
}
FeedCtrl.$inject = ['$scope', '$http', 'scroll']
