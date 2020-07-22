/*
** To browse through list of following tags and unfollow them
** Required props:
**      props.tags -> list of tags
**      props.unFollowHandler -> function(TagNameField) to call for unfollow event
*/
class FollowingTags extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            fromIndex: 0,
            perPage: 5
        };
        this.nextPage = this.nextPage.bind(this);
        this.previousPage = this.previousPage.bind(this);
    }

    render() {
        return (
            <div>
                <h2>Tags you follow:</h2>
                <Button onClick={this.previousPage}>Prev</Button>
                <Button onClick={this.nextPage}>Next</Button>
                <TagList tags = {this.props.tags}
                         fromIndex = {this.state.fromIndex}
                         perPage = {this.state.perPage}
                         actionLabel = "Unfollow"
                         onClick = {this.props.unfollowHandler}
                />
            </div>
        );
    }

    nextPage() {
        this.setState((state, props) => {
            if(state.fromIndex + state.perPage < props.tags.length) {
                return {
                    fromIndex: state.fromIndex+state.perPage
                };
            }
            else{
                return {}
            }
        });
    }

    previousPage() {
        this.setState((state, props) => {
            return {
                fromIndex: Math.max(state.fromIndex-state.perPage, 0)
            };
        });
    }
}

/*
** The entire dashboard
** Required props:
**      props.user -> user object
**      props.unfollowHandler -> function(TagNameField) to call when unfollowing a tag
*/
class Dashboard extends React.Component {
    render() {
        if(!this.props.user) { //check if no user signed in
            return <LoginButton />;
        }
        else {
            return (
                <div>
                    Hello {this.props.user.email}!
                    <br/>
                    <LogoutButton />
                    <FollowingTags
                        tags={this.props.user.tags}
                        unfollowHandler = {this.props.unfollowHandler}
                    />
                </div>
                );
        }
    }
}