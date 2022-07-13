package com.gusrylmubarok.health.backend.repositories;

import com.gusrylmubarok.health.backend.domain.User;

import java.util.List;

public interface UserDAO {
    User save(User user);
    List<User> findByEmail(String email);
    List<User> findByEmailAndPassword(String email, String password);
    void update(User user);
}
