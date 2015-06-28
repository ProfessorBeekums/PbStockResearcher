var NoteElement = React.createClass({
render: function() {
// TODO convert timestamp to date and sort
  return (
    <div className="note">
      <div className="cik">
        CIK: {this.props.cik}
      </div>
      <div className="company-name">
        Company Name: {this.props.company}
      </div>
      <div className="note-text">
        Note: {this.props.note}
      </div>
      <div className="timestamp">
        Timestamp: {this.props.timestamp}
      </div>
        <br/>
    </div>
  );
}
});
var NoteList = React.createClass({
render: function() {
  var cikFilter = this.props.cikFilter;
  var noteList = this.props.data.map(function(note) {
    if(cikFilter == '' || note.CIK == cikFilter) {
      return (
        <NoteElement key={note.NoteId} cik={note.CIK} company={note.CompanyName} note={note.Note}
        timestamp={note.Timestamp}/>
      );
    }
  });
  return (
    <div className="noteList">
      {noteList}
    </div>
  );
}
});
var NoteForm = React.createClass({
handleSubmit: function(e) {
  e.preventDefault();
  var cik = React.findDOMNode(this.refs.cik).value.trim();
  var note = React.findDOMNode(this.refs.note).value.trim();
  if (!note || !cik) {
    return;
  }
  this.props.onNoteSubmit({cik:cik, note:note});
  React.findDOMNode(this.refs.cik).value = '';
  React.findDOMNode(this.refs.note).value = '';
  return;
},
render: function() {
  return (
    <form className="noteForm" onSubmit={this.handleSubmit}>
      <input type="text" placeholder="Enter cik" ref="cik" />
      <input type="text" placeholder="Enter note text" ref="note" />
      <input type="submit" value="Post" />
    </form>
  );
}
});
var NoteFilterElement = React.createClass({
filterClicked: function(e) {
  var cik = this.props.cik;
  this.props.onFilterClicked({cik: cik});
},
render: function() {
  return (
    <p onClick={this.filterClicked}>
      {this.props.company}
    </p>
  );
}
});

var NoteFilter = React.createClass({
handleSubmit: function(e) {
  e.preventDefault();
  var cik = React.findDOMNode(this.refs.cik).value.trim();

  this.props.onNoteFilterSubmit({cik:cik});
  React.findDOMNode(this.refs.cik).value = cik;
  return;
},
handleAddFilterClick: function(e) {
  var cik = React.findDOMNode(this.refs.cik).value.trim();
  this.props.onAddFilterClick({cik: cik});
},
onFilterClicked: function(cikData) {
  this.props.onNoteFilterSubmit(cikData);
  React.findDOMNode(this.refs.cik).value = cikData.cik;
},
render: function() {
  var filterClickFunction = this.onFilterClicked;
  var filterList = this.props.filters.map(function(filter) {
    return (
      <NoteFilterElement key={filter.CompanyName} company={filter.CompanyName} cik={filter.CIK}
      onFilterClicked={filterClickFunction}/>
    );
  });
  return (
    <div className="notFilter">
      <form className="noteFilterForm" onSubmit={this.handleSubmit}>
        <input type="text" placeholder="Enter cik" ref="cik" />
        <input type="submit" value="Filter" />
      </form>
      <button type="button" onClick={this.handleAddFilterClick}>Save Filter</button>
      <div className="filterList">
        {filterList}
      </div>
    </div>
  );
}
});

var AddCompanyForm = React.createClass({
handleSubmit: function(e) {
  e.preventDefault();
  var cik = React.findDOMNode(this.refs.cik).value.trim();
  var company = React.findDOMNode(this.refs.company).value.trim();
  if (!company || !cik) {
    return;
  }
  this.props.onAddCompanySubmit({cik:cik, company:company});
  React.findDOMNode(this.refs.cik).value = '';
  React.findDOMNode(this.refs.company).value = '';
  return;
},
render: function() {
  return (
    <form className="add-company-form" onSubmit={this.handleSubmit}>
      Add a new company
      <br/>
      <input type="text" placeholder="Enter cik" ref="cik" />
      <input type="text" placeholder="Enter company name" ref="company" />
      <input type="submit" value="Post" />
    </form>
  );
}
});

var NoteFilterElement = React.createClass({
filterClicked: function(e) {
  var cik = this.props.cik;
  this.props.onFilterClicked({cik: cik});
},
render: function() {
  return (
    <p onClick={this.filterClicked}>
      {this.props.company}
    </p>
  );
}
});

var NoteBox = React.createClass({
loadNotesFromServer: function() {
  getNoteList(this.setState.bind(this))
},
handleNoteSubmit: function(noteData) {
$.ajax({
  url: 'note',
  dataType: 'json',
  type: 'POST',
  data: noteData,
  success: function(data) {
    this.loadNotesFromServer()
  }.bind(this),
  error: function(xhr, status, err) {
    console.error(this.props.url, status, err.toString());
  }.bind(this)
});
},
loadFiltersFromServer: function() {
  $.ajax({
    url: 'note-filter',
    dataType: 'json',
    cache: false,
    success: function(data) {
      this.setState({savedFilters: data});
    }.bind(this),
    error: function(xhr, status, err) {
      console.error(this.props.url, status, err.toString());
    }.bind(this)
  });
},
handleFilterSubmit: function(data) {
  this.setState({cikFilter: data.cik});
},
handleAddFilterClick: function(filter) {
$.ajax({
  url: 'note-filter',
  dataType: 'json',
  type: 'POST',
  data: filter,
  success: function(data) {
    this.loadFiltersFromServer()
  }.bind(this),
  error: function(xhr, status, err) {
    console.error(this.props.url, status, err.toString());
  }.bind(this)
});
},
handleAddCompanySubmit: function(companyData) {
$.ajax({
  url: 'company',
  dataType: 'json',
  type: 'POST',
  data: companyData,
  success: function(data) {
  }.bind(this),
  error: function(xhr, status, err) {
    console.error(this.props.url, status, err.toString());
  }.bind(this)
});
},
getInitialState: function() {
  return {data: [], cikFilter: '', savedFilters: []};
},
componentDidMount: function() {
  this.loadNotesFromServer();
  this.loadFiltersFromServer();

},
render: function() {
  return (
    <div className="noteBox">
      <div className="noteControls">
        <NoteForm onNoteSubmit={this.handleNoteSubmit}/>
        <NoteFilter filters={this.state.savedFilters} onNoteFilterSubmit={this.handleFilterSubmit}
        onAddFilterClick={this.handleAddFilterClick}/>
      </div>
      <div className="noteList">
        <NoteList data={this.state.data} cikFilter={this.state.cikFilter}/>
      </div>
      <div className="addCompanyForm">
        <AddCompanyForm onAddCompanySubmit={this.handleAddCompanySubmit}/>
      </div>
    </div>
  );
}
});

function getNoteList(successCallback) {
// TODO add ability to filter by date
// TODO will need to page (or infinite scroll) this eventually
  $.ajax({
    url: 'note',
    dataType: 'json',
    cache: false,
    success: function(data) {
      successCallback({data: data});
    },
    error: function(xhr, status, err) {
      console.error(this.props.url, status, err.toString());
    }
  });
}