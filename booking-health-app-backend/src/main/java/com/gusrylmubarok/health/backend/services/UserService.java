package com.gusrylmubarok.health.backend.services;

import com.gusrylmubarok.health.backend.domain.User;
import com.gusrylmubarok.health.backend.exceptions.UnmatchingUserCredentialsException;
import com.gusrylmubarok.health.backend.exceptions.UserNotFoundException;

public interface UserService {
    User save(User user);
    void update(User user);
    User doesUserExist(String email) throws UserNotFoundException;
    User getByEmail(String email) throws  UserNotFoundException;
    User isValidUser(String email, String password) throws UnmatchingUserCredentialsException;
}
