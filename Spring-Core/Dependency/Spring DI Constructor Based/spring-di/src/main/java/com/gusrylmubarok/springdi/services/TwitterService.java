package com.gusrylmubarok.springdi.services;

import org.springframework.stereotype.Service;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Service("TwitterService")
public class TwitterService implements MessageService{
    @Override
    public void sendMsg(String message) {
        System.out.println(message);
    }
}
