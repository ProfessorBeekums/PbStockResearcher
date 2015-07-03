var JobDash = React.createClass({
	render: function() {
		return (
		<div>
			<JobForm />
			<JobList />
		</div>
		);
	}
});

var JobList = React.createClass({
	render: function() {
		return(
			<div id="jobList">
			I should be making a request to get available jobs
			</div>
		);
	}
});

var JobForm = React.createClass({
	render: function() {
		return(
			<div id="jobForm">
				I should be a form that starts new jobs
			</div>
		);
	}
});