var JobDash = React.createClass({
	loadJobsFromServer: function() {
		getJobData(this.setState.bind(this))
	},
	onScraperSubmit: function(scraperData) {
		$.ajax({
			url: 'jobs/scraper',
			dataType: 'json',
			type: 'POST',
			data: scraperData,
			success: function(data) {
				this.loadJobsFromServer()
			}.bind(this),
			error: function(xhr, status, err) {
				console.error(this.props.url, status, err.toString());
			}.bind(this)
		});
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
			<JobForm onScraperSubmit={this.onScraperSubmit} />
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
handleSubmit: function(e) {
  e.preventDefault();
  var year = React.findDOMNode(this.refs.year).value.trim();
  var quarter = React.findDOMNode(this.refs.quarter).value.trim();
  if (!year || !quarter) {
    return;
  }
  this.props.onScraperSubmit({year:year, quarter:quarter});
  React.findDOMNode(this.refs.year).value = '';
  React.findDOMNode(this.refs.quarter).value = '';
  return;
},
render: function() {
  return (
    <form className="scraperJobsForm" onSubmit={this.handleSubmit}>
      <input type="text" placeholder="Enter year" ref="year" />
      <input type="text" placeholder="Enter quarter" ref="quarter" />
      <input type="submit" value="Post" />
    </form>
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