'use strict';

const DASHBOARD = 0;

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
                </div>);
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