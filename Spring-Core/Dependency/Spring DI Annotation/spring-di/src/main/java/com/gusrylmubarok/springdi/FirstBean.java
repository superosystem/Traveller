package com.gusrylmubarok.springdi;

import org.springframework.beans.factory.annotation.Autowired;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class FirstBean {
    @Autowired
    private SecondBean secondBean;

    @Autowired
    private ThirdBean thirdBean;

    public FirstBean() {
        System.out.println("FirstBean Initialized via Constuctor");
    }

    public void populateBeans() {
        secondBean.display();
        thirdBean.display();
    }
}
