'use strict';

const DASHBOARD = 0;

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            screen: DASHBOARD,
            user: null
        }
        this.updateUser = this.updateUser.bind(this);
        this.unfollowHandler = this.unfollowHandler.bind(this);

        this.updateUser();
    }

    // Calls API to get the user object
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

    // Calls API to unfollow a tag
    unfollowHandler(tagName) {
        console.log("Trying to unfollow"+tagName);
        $.post(`/api/tag/unfollow/${tagName}`, data => {
            if(data.result == 'success') {
                console.log(`Tag ${tagName} unfollowed.`);
            }
            this.updateUser();
        }, 'json');
    }
    render() {
        if(this.state.screen == DASHBOARD) {
            return (
                <Dashboard
                    user = {this.state.user}
                    updateUser = {this.updateUser}
                    unfollowHandler = {this.unfollowHandler}
                />
            );
        }
    }
}
ReactDOM.render(<App/>, document.getElementById('root'));