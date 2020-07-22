'use strict';

const DASHBOARD = 0;

class FollowingTags extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            fromIndex: 0,
            toIndex: this.props.perPage
        };
    }

    render() {
        return (
            <div>
                <ul>
                    {this.props.tags.slice(this.state.fromIndex, this.state.toIndex).map( tag => (
                        <li key={tag}>
                            <a href="#">{tag}</a>
                        </li>
                    ));}
                </ul>
                <button onClick={this.previousPage}>Prev</button>
                <span
                <button onClick={this.nextPage}>Next</button>
            </div>
        );
    }

    nextPage() {
        this.setState((state, props) => {
            return {
                fromIndex: state.fromIndex+props.perPage,
                toIndex: min(state.toIndex+props.perPage, props.tags.length)
            }
        });
    }

    previousPage() {

    }
}

class Dashboard extends React.Component {
    render() {
        if(!this.props.user) { //check if no user signed in
            return <a href="/login">Login</a>;
        }
        else {
            return (
                <div>
                    Hello {this.props.user.email}!
                    <br/>
                    <a href="/logout" onClick={()=>{Cookies.remove('auth-session');}}>Logout</a>
                    <FollowingTags
                        tags={this.props.user.tags}
                        perPage={5}
                    />
                </div>
                );
        }
    }
}

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            screen: DASHBOARD,
            user: null
        }
        this.updateUser();
    }
    updateUser() {
        $.getJSON('/user/get', data => {
            if(data.result == 'success') {
                this.setState({user: data.data});
            }
            else {
                this.setState({user:null});
            }
        });
    }
    render() {
        if(this.state.screen == DASHBOARD) {
            return (
                <Dashboard
                    user = {this.state.user}
                    updateUser = {this.updateUser}
                />
            );
        }
    }
}
ReactDOM.render(<App/>, document.getElementById('root'));