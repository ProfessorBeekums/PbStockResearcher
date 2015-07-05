var JobDash = React.createClass({
	loadJobsFromServer: function() {
	  getJobData(this.setState.bind(this))
	},
	getInitialState: function() {
	  return {data: []};
	},
	componentDidMount: function() {
	  this.loadJobsFromServer();

	},
	render: function() {
		return (
		<div>
			<JobForm />
			<JobList data={this.state.data} />
		</div>
		);
	}
});

var JobList = React.createClass({
	render: function() {
		var jobList = this.props.data.map(function(job) {
		    return (
		    	<tr id="job-id-{job.JobId}">
		    		<td>
		    			{job.JobType}
		    		</td>
		    		<td>
		    			{job.JobStatus}
		    		</td>
		    		<td>
		    			{job.Params}
		    		</td>
		    		<td>
		    			{job.StartTime}
		    		</td>
		    		<td>
		    			{job.EndTime}
		    		</td>
		    	</tr>
		    );
	    });
		return(
			<div id="jobList">
				<table id="job-list-table">
					<tr>
						<th>Type</th>
						<th>Status</th>
						<th>Params</th>
						<th>Start Time</th>
						<th>End Time</th>
					</tr>
					{jobList}
				</table>
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

function getJobData(successCallback) {
  $.ajax({
    url: 'jobs',
    dataType: 'json',
    cache: false,
    success: function(data) {
    	console.log(data);
      successCallback({data: data});
    },
    error: function(xhr, status, err) {
      console.error(this.props.url, status, err.toString());
    }
  });
}