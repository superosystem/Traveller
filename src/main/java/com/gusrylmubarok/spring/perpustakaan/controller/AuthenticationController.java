package com.gusrylmubarok.spring.perpustakaan.controller;

import com.gusrylmubarok.spring.perpustakaan.common.RestResult;
import com.gusrylmubarok.spring.perpustakaan.entity.User;
import com.gusrylmubarok.spring.perpustakaan.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("auth")
public class AuthenticationController extends BaseController{
    @Autowired
    private UserService service;

    @PreAuthorize("permitAll()")
    @PostMapping(value = "do-login")
    public RestResult doLogin(@RequestBody User user) {
        return service.login(user);
    }

    @PreAuthorize("permitAll()")
    @PostMapping(value = "do-register")
    public RestResult doRegister(@RequestBody User param) {
        return new RestResult(service.register(param, User.Role.ROLE_USER));
    }
}

