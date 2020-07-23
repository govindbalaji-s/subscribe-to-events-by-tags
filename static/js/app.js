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
        this.onTagDetails = this.onTagDetails.bind(this);

        this.onEventDetails = this.onEventDetails.bind(this);
        this.onEventSubscribe = this.onEventSubscribe.bind(this);
        this.onEventUnsubscribe = this.onEventUnsubscribe.bind(this);

        this.updateUser();
    }

    // Calls API to get the user object
    updateUser() {
        $.get('/user/get', data => {
            if(data.result == 'success') {
                data.data.subscribedEventsData = [];
                this.setState({user: data.data});
                //subsrcibed events ids are only there, fetch them all
                for(let i in data.data[APIUserSubscribedEventsField]) {
                    const eventid = data.data[APIUserSubscribedEventsField][i]
                    $.get(`/api/event/get/${eventid}`, data => {
                        if(data.result != 'success') {
                            console.log(`Couldn't fetch event ${eventid}`);
                        }
                        else {
                            this.setState(prevState => {
                                prevState.user.subscribedEventsData[i] = data.data;
                                return prevState;
                            });
                        }
                    }, 'json');
                }
            }
            else {
                this.setState({user:null});
            }
        });
    }

    onEventDetails(eventid) {

    }

    onEventSubscribe(eventid) {
        $.post(`/api/event/subscribe/${eventid}`, data => {
            console.log(`Subscribing ${eventid} was ${data.result}`);
            this.updateUser();
        }, 'json');
    }

    onEventUnsubscribe(eventid) {
        console.log(this);
        $.post(`/api/event/unsubscribe/${eventid}`, data => {
            console.log(`Unsubscribing ${eventid} was ${data.result}`);
            console.log(this);
            this.updateUser();
        }, 'json');
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

    onTagDetails(tagName) {

    }

    render() {
        if(this.state.screen == DASHBOARD) {
            return (
                <Dashboard
                    user = {this.state.user}
                    updateUser = {this.updateUser}
                    unfollowHandler = {this.unfollowHandler}
                    onTagDetails = {this.onTagDetails}
                    onEventDetails = {this.onEventDetails}
                    onEventSubscribe = {this.onEventSubscribe}
                    onEventUnsubscribe = {this.onEventUnsubscribe}
                />
            );
        }
    }
}
ReactDOM.render(<App/>, document.getElementById('root'));