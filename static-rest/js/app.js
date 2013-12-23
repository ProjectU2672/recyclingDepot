//(function() {
	var Todo = Backbone.Model.extend({
		defaults: {
			title: '',
			completed: false
		}
	});
	
	var TodosCollection = Backbone.Collection.extend({
		model: Todo
	});
	
	var ListView = Backbone.View.extend({
		
		template: _.template($('#item-template').html()),
		
		render: function() {
			this.$el.html(this.template(this.model.attributes));
			return this;
		}
	});
	
	var TodoView = Backbone.View.extend({
		tagName: 'li',
		
		todoTpl: _.template("An example template"),
		
		events: {
			'dblclick label':	'edit',
			'keypress .edit':	'updateOnEnter',
			'blur .edit':		'close'
		},
		
		render: function() {
			this.$el.html(this.todoTpl(this.model.toJSON()));
			this.input = this.$('.edit');
			return this;
		},
		
		edit: function() {
			console.log("EDIT");
		},
		
		close: function() {
			console.log("CLOSE");
		},
		
		updateOnEnter: function() {
			console.log("UPDATEONENTER");
		},
	});
	
	TodosView = Backbone.View.extend({
		tagName: 'ul',
		className: 'container',
		id: 'todos'
	});
//})();
