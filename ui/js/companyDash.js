var CompanyDash = React.createClass({
getInitialState: function() {
  return {data: [], cikFilter: '', savedFilters: []};
},
componentDidMount: function() {
  getNoteList(this.setState.bind(this));
},
render: function() {
// TODO need to set the filter here
  return (
    <div className="companyDash">
      <div className="noteList">
        <NoteList data={this.state.data} cikFilter={this.state.cikFilter}/>
      </div>
    </div>
  );
}
});