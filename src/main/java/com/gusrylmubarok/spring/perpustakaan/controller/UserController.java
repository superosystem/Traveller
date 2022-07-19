package com.gusrylmubarok.spring.perpustakaan.controller;

import com.gusrylmubarok.spring.perpustakaan.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("users")
public class UserController extends BaseController {
    @Autowired
    private UserService service;
}

