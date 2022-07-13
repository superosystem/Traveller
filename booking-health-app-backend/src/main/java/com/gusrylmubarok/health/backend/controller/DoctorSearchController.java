package com.gusrylmubarok.health.backend.controller;

import com.gusrylmubarok.health.backend.domain.Doctor;
import com.gusrylmubarok.health.backend.helpers.DoctorInfo;
import com.gusrylmubarok.health.backend.services.DoctorService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.ui.ModelMap;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
public class DoctorSearchController {
    final static Logger logger = LoggerFactory.getLogger(DoctorSearchController.class);

    @Autowired
    DoctorService docService;

    @GetMapping(value="/doctor/count")
    public DoctorInfo getDoctorsCount(ModelMap model) {
        int count = docService.findCount();
        return new DoctorInfo("All doctors count", count);
    }

    @RequestMapping(value="/doctor/{code}", method=RequestMethod.GET)
    public DoctorInfo getBySpecialityCode(ModelMap model, @PathVariable("code") String code) {
        List<Doctor> doctors = docService.findBySpeciality(code);
        if(doctors == null) {
            return new DoctorInfo("No Doctors found!", null);
        }
        return new DoctorInfo("Doctors found", doctors);
    }

    @GetMapping(value="/doctor", produces="application/json")
    public DoctorInfo findAll(ModelMap model) {

        List<Doctor> doctors = docService.findAll();
        if(doctors == null) {
            return new DoctorInfo("No Doctors found!", null);
        }
        return new DoctorInfo("Doctors found", doctors);
    }
}
