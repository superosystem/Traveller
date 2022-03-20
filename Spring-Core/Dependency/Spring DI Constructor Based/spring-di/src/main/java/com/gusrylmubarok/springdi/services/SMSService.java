package com.gusrylmubarok.springdi.services;

import org.springframework.stereotype.Service;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Service("SMSService")
public class SMSService implements MessageService{
    @Override
    public void sendMsg(String message) {
        System.out.println(message);
    }
}
