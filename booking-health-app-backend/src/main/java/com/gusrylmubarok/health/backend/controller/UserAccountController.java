package com.gusrylmubarok.health.backend.controller;

import com.gusrylmubarok.health.backend.domain.User;
import com.gusrylmubarok.health.backend.exceptions.UserNotFoundException;
import com.gusrylmubarok.health.backend.helpers.ExecutionStatus;
import com.gusrylmubarok.health.backend.services.DoctorService;
import com.gusrylmubarok.health.backend.services.UserService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.ui.ModelMap;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.ModelAndView;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.servlet.http.HttpSession;

import java.util.Collections;
import java.util.Map;

@RestController
@RequestMapping("/account/*")
public class UserAccountController {

    final static Logger logger = LoggerFactory.getLogger(UserAccountController.class);

    UserService userService;
    DoctorService docService;

    @Autowired
    public UserAccountController(UserService userService, DoctorService docService) {
        this.userService = userService;
        this.docService = docService;
    }

    @GetMapping(value="/token")
    public Map<String, String> getToken(HttpSession session) {
        return Collections.singletonMap("token", session.getId());
    }

    @PostMapping(value="/signup")
    public ExecutionStatus processSignup(ModelMap model, @RequestBody User reqUser) {
        User user = null;
        try {
            user = userService.doesUserExist(reqUser.getEmail());
        } catch (UserNotFoundException e) {
        }
        if(user != null) {
            return new ExecutionStatus("USER_ACCOUNT_EXISTS", "User account with same email address exists. Please try again!");
        }
        user = new User();
        user.setEmail(reqUser.getEmail());
        user.setPassword(reqUser.getPassword());
        user.setFirstName(reqUser.getFirstname());
        user.setLastname(reqUser.getLastname());
        user.setContactNumber(reqUser.getContactNumber());
        user.setAlternateContactNumber(reqUser.getAlternateContactNumber());
        user.setCityCode(reqUser.getCityCode());
        user.setStateCode(reqUser.getStateCode());
        user.setCountryCode(reqUser.getCountryCode());
        user.setAge(reqUser.getAge());
        user.setGender(reqUser.getGender());
        user.setRole(reqUser.getRole());
        User persistedUser = userService.save(user);
        //
        // Add entry in doctor table if user is a doctor
        //
        docService.addDoctor(user);

        return new ExecutionStatus("USER_ACCOUNT_CREATED", "User account successfully created");
    }

    @PostMapping(value="/user/update")
    public ExecutionStatus updateUser(ModelMap model, @RequestBody User reqUser) {
        User user = new User();
        user.setId(reqUser.getId());
        user.setFirstName(reqUser.getFirstname());
        user.setLastname(reqUser.getLastname());
        user.setContactNumber(reqUser.getContactNumber());
        user.setAlternateContactNumber(reqUser.getAlternateContactNumber());
        user.setCityCode(reqUser.getCityCode());
        user.setStateCode(reqUser.getStateCode());
        user.setCountryCode(reqUser.getCountryCode());
        user.setAge(reqUser.getAge());
        user.setGender(reqUser.getGender());
        userService.update(user);
        return new ExecutionStatus("USER_ACCOUNT_UPDATED", "User account successfully updated");
    }

    @PostMapping(value="/update", produces="application/json")
    public ModelAndView updateProfile(ModelMap model, @RequestParam("firstName") String firstName,
                                      @RequestParam("lastName") String lastName, @RequestParam("address") String address,
                                      @RequestParam("contact_number") String contactNumber) {
        return new ModelAndView("update", model);
    }

    @PostMapping(value="/forgotpassword/process", produces="application/json")
    public ModelAndView processForgotPassword(ModelMap model, @RequestParam("emailaddress") String email) {

        User user = null;
        try {
            user = userService.doesUserExist(email);
        } catch (UserNotFoundException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
        if(user != null) {

        }
        model.addAttribute("message", "An email notification is sent to the registered email address.");
        return new ModelAndView("forgotpassword", model);
    }
}